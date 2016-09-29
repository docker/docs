open Lwt.Infix

let src =
  let src = Logs.Src.create "volumes" ~doc:"map volumes to the VM" in
  Logs.Src.set_level src (Some Logs.Info);
  src

module Log = (val Logs.src_log src : Logs.LOG)

type key = string

type t = {
  container_id: string;
  host_paths: string list;
}

let get_key t = t.container_id

module Map = Map.Make(String)

type 'a tree = {
  node : 'a;
  children : 'a tree Map.t;
}

(* When event streams are coalesced hierarchically, descendent streams
   become subscriptions while ancestor streams flow. If all
   ancestor streams terminate before a descendent stream, the subscription
   re-hydrates into a first-class FS event stream. *)
type watch_action = {
  wpath : string;
  write : Cstruct.t -> int Lwt.t;
  since : Fsevents.EventId.t;
  drain : Fsevents.EventId.t option tree;
}

type action = Watch_action of watch_action | Subscribe_action of string

type watch_state = Starting | Watching | Stopping | Stopped

type watch = {
  runloop : Cf.RunLoop.t;
  watcher : Fsevents_lwt.t;
  send    : Cstruct.t -> int Lwt.t;
  state   : watch_state ref;
}

type sub = Watch of watch | Subscribe

type subscription = {
  path        : string;
  subscribers : string list;
  mutable sub : sub option;
}

type unwatch_action =
  | Live of subscription tree
  | Dead

type syscall =
  | Ping
  | Rmdir
  | Unlink
  | Mkdir
  | Symlink
  | Truncate
  | Chmod_dir
  | Chmod
  | Mknod_reg

let subscription_empty = { path = ""; subscribers = []; sub = None; }

type context = {
  event_fd : Lwt_unix.file_descr;
  mounts : Mount.t list;
  ctl : Control.t;
}

type state = {
  mutable watches : subscription tree;
  mutable containers : t Map.t;
}

let state = {
  watches = { node = subscription_empty; children = Map.empty };
  containers = Map.empty;
}

let syscall_stream, syscall_push' = Lwt_stream.create ()

let syscall_push = ref (fun x -> syscall_push' x; Lwt.return_unit)

let string_of_id = function
  | Fsevents.EventId.Now -> "now"
  | Fsevents.EventId.Since t -> Unsigned.UInt64.to_string t

let string_of_syscall = function
  | Ping      -> "ping"
  | Rmdir     -> "rmdir"
  | Unlink    -> "unlink"
  | Mkdir     -> "mkdir"
  | Symlink   -> "symlink"
  | Truncate  -> "truncate"
  | Chmod_dir -> "chmod_dir"
  | Chmod     -> "chmod"
  | Mknod_reg -> "mknod_reg"

let to_string t = Printf.sprintf "%s" t.container_id

let description_of_format = "container_id:path1:path2:path3"

let run_loop_mode = Cf.RunLoop.Mode.Default
let create_flags = Fsevents.CreateFlags.({
  use_cf_types = false;
  no_defer = true;
  watch_root = true;
  ignore_self = true;
  file_events = true;
  mark_self = false;
})

let enum_of_syscall = function
  | Ping -> 128
  | Rmdir -> 0
  | Unlink -> 1
  | Mkdir -> 2
  | Symlink -> 3
  | Truncate -> 4
  | Chmod_dir
  | Chmod -> 5 (* Actuation is identical but distinguished for placation *)
  | Mknod_reg -> 6

(* 2 bytes for total length
   2 bytes for path length + NUL (x)
   x bytes for path
   1 byte  for syscall
*)
let write_syscall path write syscall =
  let plen = String.length path + 1 in
  let len = 4 + 2 + plen + 1 in
  let b = Cstruct.create len in
  Cstruct.LE.set_uint32 b 0 (Int32.of_int len);
  Cstruct.BE.set_uint16 b 4 plen;
  Cstruct.blit_from_string path 0 b 6 (plen - 1);
  Cstruct.set_uint8 b (plen + (6 - 1)) 0;
  Cstruct.set_uint8 b (plen + 6) (enum_of_syscall syscall);
  (*Log.err (fun f -> f "event message: %S" (Cstruct.to_string b));*)
  Log.debug (fun f ->
    f "sending syscall %s for %s" (string_of_syscall syscall) path);
  !syscall_push (Some (path, syscall))
  >>= fun () ->
  write b
  >>= fun written ->
  (if written <> len
   then Log.err (fun f -> f "INCOMPLETE EVENT WRITE: %d of %d" written len)
  );
  Lwt.return_unit

let log_event event =
  let open Fsevents_lwt in
  Log.debug (fun f -> f "FSEvent %s %s %s"
                (Unsigned.UInt64.to_string
       (Fsevents.EventId.to_uint64 event.id))
    event.path
    (Fsevents.EventFlags.to_string_one_line event.flags))

let remove_syscall event = Fsevents.EventFlags.(match event.item_type with
  | Some Dir -> Rmdir
  | None | Some (File | Symlink | Hardlink) -> Unlink
)

let create_syscall event = Fsevents.(match event.EventFlags.item_type with
  | Some EventFlags.Dir -> Mkdir
  | Some EventFlags.Symlink -> Symlink
  | None
  | Some (EventFlags.File | EventFlags.Hardlink) -> Mknod_reg
)

let modified_syscall event = Fsevents.EventFlags.(match event.item_type with
  | Some Dir -> Mkdir
  | None | Some (File | Symlink | Hardlink) -> Truncate
)

let attrib_syscall event = Fsevents.(match event.EventFlags.item_type with
  | Some EventFlags.Dir -> Chmod_dir
  | None
  | Some ( EventFlags.File
         | EventFlags.Symlink
         | EventFlags.Hardlink) -> Chmod
)

(* TODO: check all the cases here! *)
let syscalls_for_root_changed event =
  Lwt.catch (fun () ->
    Lwt_unix.stat event.Fsevents_lwt.path
    >>= fun _stat ->
    Lwt.return []
  ) (function
    | Unix.Unix_error (Unix.ENOENT, "stat", _) ->
      (* TODO: THIS COULD BE A FILE!!! *)
      Lwt.return [Rmdir]
    | exn -> Lwt.fail exn
  )

let rec syscalls_of_event event syscalls flags =
  let open Fsevents_lwt in
  let open Fsevents.EventFlags in
  match flags with
  | { item_created = true; item_removed = true } ->
    Lwt.catch
      (fun () ->
         Lwt_unix.lstat event.path
         >>= fun _stat ->
         syscalls_of_event event
           ((remove_syscall flags)::syscalls)
           { flags with item_removed = false }
      )
      (function
        | Unix.Unix_error (Unix.ENOENT, _, _) ->
          syscalls_of_event event
            ((create_syscall flags)::syscalls)
            { flags with item_created = false }
        | exn ->
          Log.err (fun f ->
            f "UNEXPECTED EXCEPTION getting mk/lstat of %s:\n%s"
              event.path (Printexc.to_string exn));
          Lwt.return (List.rev syscalls)
      )

  | { item_created = true } ->
    syscalls_of_event event
      ((create_syscall flags)::syscalls)
      { flags with item_created = false }

  | { item_inode_meta_mod = true }
  | { item_finder_info_mod = true }
  | { item_change_owner = true }
  | { item_xattr_mod = true } ->
    syscalls_of_event event
      ((attrib_syscall flags)::syscalls)
      { flags with
        item_inode_meta_mod = false;
        item_finder_info_mod = false;
        item_change_owner = false;
        item_xattr_mod = false;
      }

  | { item_modified = true; item_type = Some Dir } ->
    log_event event;
    Lwt.return (List.rev syscalls)

  | { item_modified = true } ->
    syscalls_of_event event
      ((modified_syscall flags)::syscalls)
      { flags with item_modified = false }

  | { item_removed = true } ->
    syscalls_of_event event
      ((remove_syscall flags)::syscalls)
      { flags with item_removed = false }

  | { item_renamed = true } ->
    Lwt.catch (fun () ->
      Lwt_unix.lstat event.path
      >>= fun _stat ->
      syscalls_of_event event
        ((modified_syscall flags)::syscalls)
        { flags with item_renamed = false }
    ) (function
      | Unix.Unix_error (Unix.ENOENT, _, _) ->
        syscalls_of_event event
          ((remove_syscall flags)::syscalls)
          { flags with item_renamed = false }
      | exn ->
        Log.err (fun f ->
          f "UNEXPECTED EXCEPTION getting mv/lstat of %s:\n%s"
            event.path (Printexc.to_string exn));
        Lwt.return (List.rev syscalls)
    )

  | { root_changed = true } ->
    syscalls_for_root_changed event
    >>= fun generated_syscalls ->
    syscalls_of_event event
      (generated_syscalls@syscalls)
      { flags with root_changed = false }

  | { item_type = Some _ } ->
    (* WIPE THE ITEM_TYPE IF WE DON'T KNOW HOW TO HANDLE ANY MORE FLAGS *)
    syscalls_of_event event syscalls { flags with item_type = None }
  | flags when flags = zero ->
    Lwt.return (List.rev syscalls)
  | _unhandled_flags ->
    log_event event;
    Lwt.return (List.rev syscalls)

let trigger write event =
  log_event event;
  syscalls_of_event event [] event.Fsevents_lwt.flags
  >>= fun syscalls ->
  Lwt_list.iter_s (write_syscall event.Fsevents_lwt.path write) syscalls

let rec get_threshold drain = function
  | [] -> drain.node
  | next :: rest -> match Map.find next drain.children with
    | exception Not_found -> None
    | { node = None } -> get_threshold drain rest
    | { node } -> node

let drain_is_empty = function
  | { node = None; children; } when Map.is_empty children -> true
  | _ -> false

let should_drain drain prefix_size event =
  if drain_is_empty drain then true else
    let open Fsevents_lwt in
    let path_length = String.length event.path in
    let path = String.sub event.path prefix_size (path_length - prefix_size) in
    let segments = match Stringext.split ~on:'/' path with
      | "" :: rest -> rest
      | xs -> xs
    in
    match get_threshold drain segments with
    | Some id -> Fsevents.EventId.compare event.id id > 0
    | None    -> true

let rec filter_drain id drain = {
  node = (match drain.node with
    | None -> None
    | Some Fsevents.EventId.Now ->
      None (* This un-normalization indicates a bug elsewhere. *)
    | Some threshold when Fsevents.EventId.compare id threshold > 0 -> None
    | Some threshold -> Some threshold
  );
  children = filter_children id drain.children;
}
and filter_children id = Map.filter (fun _ drain ->
  let drain = filter_drain id drain in
  not (drain_is_empty drain)
)

let rec continue_watch write watcher =
  let stream = Fsevents_lwt.stream watcher in
  Lwt_stream.get stream
  >>= function
  | None -> Log.debug (fun f -> f "Stream ended (continue)"); Lwt.return_unit
  | Some event ->
    trigger write event
    >>= fun () ->
    continue_watch write watcher

let rec start_watch watcher write drain prefix_size state () =
  let stream = Fsevents_lwt.stream watcher in
  Lwt_stream.get stream
  >>= function
  | None -> Log.debug (fun f -> f "Stream ended (start)"); Lwt.return_unit
  | Some event ->
    Log.debug (fun f -> f "start_watch: got event from stream");
    let drain = filter_drain event.id drain in
    (if should_drain drain prefix_size event
     then trigger write event
     else Lwt.return_unit)
    >>= fun () ->
    if event.Fsevents_lwt.flags.Fsevents.EventFlags.history_done
       || drain_is_empty drain
    then begin
      state := Watching;
      continue_watch write watcher
    end
    else start_watch watcher write drain prefix_size state ()

(* TODO: what happens when watching an ancestor of a VM bind mount? *)
(* TODO: what happens when watched tree contains mounts? *)
(* TODO: watching through symlinks results in RootChanged a lot but
   nothing else. We need to watch each intermediate symlink and take
   appropriate action when any of them change. *)
(* TODO: Check watch root changes esp racing with coalescion/withering. *)
let watch { since; wpath; drain; write } =
  Log.debug (fun f -> f "Volume: watching %s since %s" wpath (string_of_id since));
  let watcher = Fsevents_lwt.create ~since 0. create_flags [wpath] in
  let state = ref Starting in
  let prefix_size = String.length wpath in
  Lwt.async (start_watch watcher write drain prefix_size state);
  Cf_lwt.RunLoop.run_thread (fun runloop ->
    Fsevents_lwt.schedule_with_run_loop watcher runloop run_loop_mode;
    if not (Fsevents_lwt.start watcher)
    then Log.err (fun f -> f "Failed to start FSEvents stream for %s" wpath)
    else Log.debug (fun f -> f "Volume: started watcher for %s" wpath)
  ) >>= fun runloop ->
  Lwt.return { runloop; watcher; send = write; state; }

let attach container_id action sub =
  let sub = {
    sub with subscribers = container_id :: sub.subscribers;
  } in
  match action with
  | Watch_action waction ->
    watch waction
    >>= fun watch ->
    Lwt.return { sub with sub = Some (Watch watch); path = waction.wpath }
  | Subscribe_action path ->
    match sub.sub with
    | None | Some Subscribe ->
      Lwt.return { sub with sub = Some Subscribe; path; }
    | Some (Watch _) -> Lwt.return sub

let rec build_tree container_id action = function
  | [] ->
    attach container_id action subscription_empty
    >>= fun node ->
    Lwt.return { children = Map.empty; node; }
  | next :: rest ->
    build_tree container_id action rest
    >>= fun tree ->
    let children = Map.singleton next tree in
    Lwt.return { children; node = subscription_empty }

let watched = function
  | { node = { sub = Some _ } } -> true
  | { node = { sub = None   } } -> false

let tree_action tree = function
  | Subscribe_action path -> Subscribe_action path
  | Watch_action { wpath } when watched tree -> Subscribe_action wpath
  | Watch_action _waction as action -> action

let rec waction_of_root root waction =
  Map.fold (fun entry tree waction ->
    waction >>= fun waction -> match tree with
    | { node = { sub = Some (Watch { runloop; watcher; state; }); path } } ->
      state := Stopping;
      Fsevents_lwt.flush watcher
      >>= fun () ->
      Fsevents_lwt.stop watcher;
      state := Stopped;
      let event_id = Fsevents_lwt.get_latest_event_id watcher in
      Fsevents_lwt.invalidate watcher;
      Fsevents_lwt.release watcher;
      (* TODO: ugggh why is this side effect here? *)
      tree.node.sub <- Some Subscribe;
      Cf.RunLoop.stop runloop;
      let since = Fsevents.EventId.min waction.since event_id in
      Log.debug (fun f -> f "drainstop %s %s" path (string_of_id event_id));
      let node = Fsevents.EventId.(match since with
        | Now -> None
        | Since _ -> Some since
      ) in
      let node = { node; children = Map.empty; } in
      let drain = {
        waction.drain
        with children = Map.add entry node waction.drain.children
      } in
      Lwt.return { waction with drain; since; }
    | { node = { sub = (Some Subscribe | None) }; children } ->
      waction_of_root children waction
      >>= fun children_waction ->
      let drain = {
        waction.drain with
        children = Map.add entry children_waction.drain waction.drain.children
      } in
      let since = Fsevents.EventId.min waction.since children_waction.since in
      Lwt.return { waction with drain; since }
  ) root (Lwt.return waction)

let watch_action root = function
  | Subscribe_action path -> Lwt.return (Subscribe_action path)
  | Watch_action waction ->
    waction_of_root root waction
    >>= fun waction ->
    Lwt.return (Watch_action waction)

let rec watch_path tree container_id action = function
  | "" :: segments -> watch_path tree container_id action segments
  | [] ->
    watch_action tree.children action
    >>= fun action ->
    attach container_id action tree.node
    >>= fun node ->
    Lwt.return { tree with node }
  | next :: rest ->
    begin match Map.find next tree.children with
      | exception Not_found -> build_tree container_id action rest
      | branch ->
        let action = tree_action branch action in
        watch_path branch container_id action rest
    end
    >>= fun branch ->
    Lwt.return { tree with children = Map.add next branch tree.children }

let add_watch inotify_fd t =
  state.containers <- Map.add t.container_id t state.containers;
  Lwt_list.fold_left_s (fun tree wpath ->
    (* TODO: consolidate watched trees *)
    Log.debug (fun f -> f "Adding watch for %s" wpath);
    let segments = Stringext.split ~on:'/' wpath in
    let since = Fsevents.EventId.Now in
    let drain = { node = None; children = Map.empty; } in
    let write cstruct =
      let len = Cstruct.len cstruct in
      let buf = Bytes.create len in
      Cstruct.blit_to_bytes cstruct 0 buf 0 len;
      Lwt_unix.write inotify_fd buf 0 len
    in
    let waction = Watch_action { wpath; since; drain; write; } in
    watch_path tree t.container_id waction segments
  ) state.watches t.host_paths
  >>= fun watches ->
  state.watches <- watches;
  Lwt.return_unit

let rec restart write since children = Map.fold (fun key tree map_lwt ->
  map_lwt >>= fun map ->
  match tree.node.sub with
  | None ->
    restart write since tree.children
    >>= fun children ->
    Lwt.return (Map.add key { tree with children } map)
  | Some Subscribe ->
    let drain = { node = None; children = Map.empty; } in
    Log.debug (fun f -> f "restarting watch of %s" tree.node.path);
    watch { since; wpath=tree.node.path; drain; write; }
    >>= fun w ->
    tree.node.sub <- Some (Watch w);
    Lwt.return map
  | Some (Watch _) ->
    Log.err (fun f -> f "WATCH FRONTIER INVARIANT VIOLATED");
    Lwt.return map
) children (Lwt.return children)

let wither write since tree =
  if Map.is_empty tree.children then Lwt.return Dead
  else
    restart write since tree.children
    >>= fun children ->
    Lwt.return (Live { tree with children })

let detach container_id tree =
  (if not (List.mem container_id tree.node.subscribers)
   then Log.err (fun f -> f "WATCH SUBSCRIBER REMOVAL INVARIANT VIOLATED")
  );
  match tree with
  | { node = { subscribers = _ :: _ :: _ } as sub } ->
    let subscribers = List.filter ((<>) container_id) tree.node.subscribers in
    let node = { sub with subscribers; } in
    Lwt.return (Live { tree with node })
  | { node = { sub = (Some Subscribe | None) } } ->
    let liveness =
      if Map.is_empty tree.children then Dead
      else Live { tree with node = subscription_empty }
    in
    Lwt.return liveness
  | { node = { sub = Some (Watch watch); path; } as sub } ->
    watch.state := Stopping;
    Log.debug (fun f -> f "about to flush watcher for %s" path);
    Fsevents_lwt.flush watch.watcher
    >>= fun () ->
    Fsevents_lwt.stop watch.watcher;
    watch.state := Stopped;
    let since = Fsevents_lwt.get_latest_event_id watch.watcher in
    Fsevents_lwt.invalidate watch.watcher;
    Fsevents_lwt.release watch.watcher;
    Log.debug (fun f -> f "stopped watching %s" path);
    sub.sub <- None; (* TODO: ugggh why is this side effect here? *)
    Cf.RunLoop.stop watch.runloop;
    wither watch.send since tree

let rec unwatch_path tree container_id = function
  | [] -> detach container_id tree
  | "" :: segments -> unwatch_path tree container_id segments
  | next :: rest ->
    begin match Map.find next tree.children with
      | exception Not_found ->
        Log.err (fun f -> f "WATCH TREE REMOVAL EXISTENCE INVARIANT VIOLATED");
        Lwt.return tree.children
      | branch ->
        unwatch_path branch container_id rest
        >>= function
        | Dead -> Lwt.return (Map.remove next tree.children)
        | Live branch -> Lwt.return (Map.add next branch tree.children)
    end
    >>= fun children ->
    if Map.is_empty children && tree.node.sub = None
    then Lwt.return Dead
    else Lwt.return (Live { tree with children })

let remove_watch t =
  state.containers <- Map.remove t.container_id state.containers;
  Lwt_list.fold_left_s (fun tree path ->
    (* TODO: sort watch removals *)
    Log.debug (fun f -> f "Removing watch for %s" path);
    let segments = Stringext.split ~on:'/' path in
    unwatch_path tree t.container_id segments
    >>= function
    | Dead ->
      Log.debug (fun f -> f "removed last watch");
      Lwt.return { node = subscription_empty; children = Map.empty }
    | Live tree -> Lwt.return tree
  ) state.watches t.host_paths
  >>= fun watches ->
  state.watches <- watches;
  Lwt.return_unit

let partition_shared mounts = List.partition (fun path ->
  List.exists (fun mount ->
    let mount_path = Mount.mount_path mount in
    let mount_len = String.length mount_path in
    path = mount_path || (
      String.length path > mount_len &&
      String.sub path 0 (mount_len + 1) = (mount_path ^ "/")
    )
  ) mounts
)

let ns_doc_url = "https://docs.docker.com/docker-for-mac/osxfs/#namespaces"

let start ctxt t =
  Log.info (fun f ->
    f "Volume.start %s (paths = [%s])"
      (to_string t) (String.concat ", " t.host_paths)
  );
  Active_list.Var.read ctxt
  >>= fun { event_fd; mounts; ctl; } ->
  (* Filter out named volumes *)
  let paths = List.filter (fun s -> String.get s 0 = '/') t.host_paths in
  let host_paths, vm_paths = partition_shared mounts paths in
  Control.partition_suitable_mounts ctl vm_paths
  >>= function
  | [], _ ->
    Lwt.catch
      (fun () ->
         Lwt_list.map_p (fun p ->
           Docker_realname.resolve p >>= fun r -> Lwt.return (p, r)
         ) host_paths
         >>= fun host_paths ->
         Log.info (fun f ->
           let resolutions =
             List.map (fun (p, r) -> p ^ " -> " ^ r) host_paths
           in
           f "Volume.start %s (watches [%s])"
             (to_string t) (String.concat ", " resolutions)
         );
         let host_paths =
           List.sort String.compare (List.rev_map snd host_paths)
         in
         let ft = { t with host_paths } in
         add_watch event_fd ft
         >>= fun () ->
         Lwt.return (Result.Ok ft)
      )
      (function
        | Docker_realname.Error msg ->
          Lwt.return (Result.Error (`Msg msg))
        | exn -> Lwt.fail exn
      )
  | [missing_or_empty], _ ->
    (* The trailing dot in these messages is necessary to get the Docker
       client to somewhat reasonably return the terminal. Without the
       trailing dot, the terminal is returned without a trailing
       newline or carriage return and so the next prompt is printed
       indented which causes annoying behavior like garbled history in
       some shells. *)
    let msg =
      Printf.sprintf "\r\nThe path %s\r\nis not shared from OS X \
                      and is not known to Docker.\r\n\
                      You can configure shared paths from \
                      Docker -> Preferences... -> File Sharing.\r\n\
                      See %s for more info.\r\n."
        missing_or_empty ns_doc_url
    in
    Lwt.return (Result.Error (`Msg msg))
  | missing_or_empty, _ ->
    let msg =
      Printf.sprintf "\r\nThe paths %s\r\nare not shared from OS X \
                      and are not known to Docker.\r\n\
                      You can configure shared paths from \
                      Docker -> Preferences... -> File Sharing.\r\n\
                      See %s for more info.\r\n."
        (String.concat " and " missing_or_empty) ns_doc_url
    in
    Lwt.return (Result.Error (`Msg msg))

let stop t =
  Log.info (fun f -> f "Volume.stop %s (paths = [%s])" (to_string t) (String.concat ", " t.host_paths));
  remove_watch t
  >>= fun () ->
  Lwt.return ()

let of_string x = match Stringext.split ~on:':' x with
  | container_id :: host_paths ->
    Result.Ok { container_id; host_paths }
  | _ ->
    Result.Error (`Msg ("Failed to parse request, expected " ^ description_of_format))
