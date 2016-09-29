
module Log = Log.Log

module PathSet = Set.Make(String)

open Lwt.Infix

let error_ret_code = 0
let strings_ret_code = 1

let does_overlap path_a path_b =
  let path_a = Astring.String.trim ~drop:((=) '/') path_a in
  let path_b = Astring.String.trim ~drop:((=) '/') path_b in
  path_a = path_b ||
  Astring.String.is_prefix ~affix:(path_a^"/") path_b ||
  Astring.String.is_prefix ~affix:(path_b^"/") path_a

(* given path sets A and B, compute all pairs of B paths that overlap
   an A path *)
let path_set_overlaps a b =
  List.fold_left (fun pairs path_b ->
    let overlaps_a =
      PathSet.filter (fun path_a -> does_overlap path_a path_b) a
    in
    (List.map (fun o_a -> path_b, o_a) (PathSet.elements overlaps_a))@pairs
  ) [] (PathSet.elements b)

let unique_overlaps list =
  List.sort_uniq compare (List.fold_left (fun list (a, b) ->
    match String.compare a b with
    | 0 -> list
    | x when x < 0 -> (b, a) :: list
    | _ -> (a, b) :: list
  ) [] list)

(* check exists on host, check with VM *)
let check_mount_paths ctl mounts paths =
  let mounts = PathSet.of_list mounts in
  let paths = PathSet.of_list (List.map Bytes.to_string paths) in
  let news = PathSet.diff paths mounts in
  let olds = PathSet.diff mounts paths in
  match unique_overlaps (path_set_overlaps paths news) with
  | first :: rest ->
    let msg =
      if rest = []
      then
        Printf.sprintf "The export path %s overlaps with the export path %s."
          (fst first) (snd first)
      else
        let overlaps = unique_overlaps (first :: rest) in
        Printf.sprintf "Some export paths overlap:\n%s"
          (String.concat "\n"
             (List.map (fun (a, b) -> a ^ " overlaps " ^ b) overlaps))
    in
    Lwt.return (Result.Error msg)
  | [] ->
    Lwt_list.filter_p
      (fun path -> Lwt_unix.file_exists path >|= (not))
      (PathSet.elements news)
    >>= function
    | first :: rest ->
      let msg =
        if rest = []
        then Printf.sprintf "The export path %s does not exist on OS X." first
        else Printf.sprintf "The export paths %s and %s do not exist on OS X."
            (String.concat ", " rest) first
      in
      Lwt.return (Result.Error msg)
    | [] ->
      let over, under, checkable =
        PathSet.fold (fun new_export (over, under, checkable) ->
          match path_set_overlaps olds (PathSet.singleton new_export) with
          | [] -> (over, under, PathSet.add new_export checkable)
          | (_, old_export)::_ -> (* assumes disjoint invariant on olds *)
            if String.length old_export < String.length new_export
            then (over, PathSet.add new_export under, checkable)
            else ((new_export, old_export)::over, under, checkable)
        ) news ([], PathSet.empty, PathSet.empty)
      in
      match over with
      | (_ :: _) ->
        (* TODO: fix this case so two steps are not required *)
        let msg = String.concat "\n\n" (List.map (fun (new_path, old_path) ->
          Printf.sprintf
            "The path %s cannot be exported when %s is exported.\n%s"
            new_path old_path
            ("Please unexport "^old_path^" before exporting "^new_path^".")
        ) over)
        in
        Lwt.return (Result.Error msg)
      | [] ->
        Control.partition_suitable_exports ctl (PathSet.elements checkable)
        >>= function
        | _, bad_first :: bad_rest ->
          let see_doc = "Please see the documentation for more information." in
          let msg =
            if bad_rest = []
            then Printf.sprintf
                "The path %s is reserved by Docker.\n%s" bad_first see_doc
            else Printf.sprintf
                "The paths %s and %s are reserved by Docker.\n%s"
                (String.concat ", " bad_rest) bad_first see_doc
          in
          Lwt.return (Result.Error msg)
        | ok_news, [] ->
          let paths = List.sort String.compare (PathSet.elements paths) in
          let new_set = PathSet.of_list ok_news in
          let ok = 0 = PathSet.compare new_set checkable in
          if ok
          then Lwt.return (Result.Ok paths)
          else
            let news = List.sort String.compare ok_news in
            let msg = Printf.sprintf
                "An error occurred processing the export paths:\n%s\n\n%s\n%s"
                (String.concat "\n" paths)
                "They were returned as:"
                (String.concat "\n" news)
            in
            Lwt.return (Result.Error msg)

let export_path_syntax_errors = List.fold_left (fun errors path ->
  let path = Bytes.to_string path in
  if String.get path 0 <> '/'
  then (Printf.sprintf "%s is not an absolute path." path)::errors
  else match Astring.String.find_sub ~sub:"//" path with
    | Some _ -> (Printf.sprintf "%s contains double slashes." path)::errors
    | None ->
      let segments = Astring.String.cuts ~sep:"/" path in
      if List.exists ((=) ".") segments
      then (Printf.sprintf "%s contains a '.' path segment." path)::errors
      else if List.exists ((=) "..") segments
      then (Printf.sprintf "%s contains a '..' path segment." path)::errors
      else errors
) []

let handle_check_mount_paths ctl mounts body =
  let rec collect_strings acc off =
    try
      let noff = Bytes.index_from body off '\000' in
      let s = Bytes.sub body off (noff - off) in
      collect_strings (s::acc) (noff + 1)
    with Not_found -> List.rev acc
  in
  let paths = collect_strings [] 0 in
  (match export_path_syntax_errors paths with
   | (_ :: _) as errors ->
     let msg = Printf.sprintf "Some export paths had syntax errors:\n%s"
         (String.concat "\n" errors)
     in
     Lwt.return (Result.Error msg)
   | [] ->
     check_mount_paths ctl mounts paths
  ) >>= function
  | Result.Error error ->
    let len = String.length error in
    let msg_len = 6 + len + 1 in
    let cs = Cstruct.create msg_len in
    Cstruct.LE.set_uint32 cs 0 (Int32.of_int msg_len);
    Cstruct.LE.set_uint16 cs 4 error_ret_code;
    Cstruct.blit_from_string error 0 cs 6 len;
    Cstruct.set_uint8 cs (msg_len - 1) 0;
    Lwt.return (Cstruct.copy cs 0 msg_len)
  | Result.Ok paths ->
    let len = List.fold_left (fun a s -> a + 1 + String.length s) 0 paths in
    let msg_len = 6 + len in
    let cs = Cstruct.create msg_len in
    Cstruct.LE.set_uint32 cs 0 (Int32.of_int msg_len);
    Cstruct.LE.set_uint16 cs 4 strings_ret_code;
    ignore (List.fold_left (fun off s ->
      let len = String.length s in
      Cstruct.blit_from_string s 0 cs off len;
      Cstruct.set_uint8 cs (off + len) 0;
      off + len + 1
    ) 6 paths);
    Lwt.return (Cstruct.copy cs 0 msg_len)

(* An Lwt thread which receives commands from the socket and
   dispatches them appropriately. *)
let serve_forever ctl mounts control_sock =
  let buf = Bytes.make 6 '\000' in

  let rec client_loop control_fd =
    Lwt_unix.read control_fd buf 0 6
    >>= function
    | 0 ->
      Lwt.fail (Failure "osxfs_control: eof reading header from socket")
    | 6 ->
      let header = Cstruct.of_string (Bytes.to_string buf) in
      let len = Int32.to_int (Cstruct.LE.get_uint32 header 0) - 6 in
      let op = Cstruct.LE.get_uint16 header 4 in
      if len = 0
      then handle_op control_fd Bytes.empty op
      else
        let body = Bytes.make len '\000' in
        Lwt_unix.read control_fd body 0 len
        >>= (function
          | 0 ->
            Lwt.fail (Failure "osxfs_control: eof reading body from socket")
          | x when x = len -> handle_op control_fd body op
          | _ ->
            Lwt.fail
              (Failure "osxfs_control: partial read of control packet body")
        )
    | _ ->
      Lwt.fail (Failure "osxfs_control: partial read of control packet header")
  and handle_op control_fd body = function
    | 0 -> begin
        let mounts = List.map Mount.mount_path mounts in
        handle_check_mount_paths ctl mounts body
        >>= fun buf ->
        let len = String.length buf in
        Lwt_unix.write control_fd (Bytes.of_string buf) 0 len
        >>= function
        | 0 ->
          Lwt.fail
            (Failure "osxfs_control: eof writing response to socket")
        | x when x = len -> client_loop control_fd
        | _ ->
          Lwt.fail (Failure "osxfs_control: partial write of response")
      end
    | _ ->
      Lwt.fail (Failure "osxfs_control: unknown control packet type")
  in
  let rec server_loop () =
    Lwt_unix.accept control_sock
    >>= fun (control_fd, _sockaddr) ->
    Lwt.async (fun () -> client_loop control_fd);
    server_loop ()
  in
  server_loop ()
