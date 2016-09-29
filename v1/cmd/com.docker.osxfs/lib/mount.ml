open Lwt.Infix

module LogInfo = Log.LogInfo
module Log = Log.Log

type t = {
  export_path : string;
  export_root : string;
  mount_path : string;
}

let rec get_export_root = function
  | "" -> Lwt.return ""
  | export_path ->
    Lwt.catch (fun () ->
      Lwt_unix.readlink export_path
      >>= fun real_path ->
      let path =
        if real_path.[0] = '/' then real_path
        else Filename.(concat (dirname export_path) real_path)
      in
      Lwt.return (false, path)
    ) (function
      | Unix.Unix_error (Unix.EINVAL, "readlink", _) ->
        Lwt.return (true, export_path)
      | exn -> Lwt.fail exn
    )
    >>= function
    | (false, path) -> get_export_root path
    | (true, path) -> Lwt.return path

let make ~export_path ~mount_path =
  get_export_root export_path
  >>= fun export_root ->
  Lwt.return { export_path; export_root; mount_path; }

let export_path { export_path } = export_path
let export_root { export_root } = export_root
let mount_path { mount_path } = mount_path

let segment_name { export_path } =
  let pattern = "_" in
  let with_ = "__" in
  let escaped = Stringext.replace_all export_path ~pattern ~with_ in
  Stringext.replace_all escaped ~pattern:"/" ~with_:"_"

let ( / ) = Filename.concat

let home = try
    Sys.getenv "HOME"
  with Not_found ->
    Log.err (fun f ->
      f "Could not retrieve %s environment variable. %s is required by osxfs."
        "HOME" "HOME"
    );
    exit 1

let host_docker_app =
  home / "Library/Containers/com.docker.docker/Data", "/host_docker_app"

let default_mounts = [
  "/Users",   "/Users";
  "/Volumes", "/Volumes";
  "/tmp",     "/tmp";
  "/private", "/private";
  host_docker_app;
]

module Database = Client9p_unix.Make(LogInfo)

let txn_name = "osxfs-mounts"
let txn = ["branch"; "master"; "transactions"; txn_name]
let config_path = [
  "com.docker.driver.amd64-linux"; "mounts";
]

let rw = [ `Read; `Write ]
let rwx = `Execute :: rw
let txn_perm = Protocol_9p.Types.FileMode.make ~owner:rwx ~is_directory:true ()
let mounts_perm = Protocol_9p.Types.FileMode.make ~owner:rw ()

let read_key db path =
  (* TODO: handle configs larger than 4KiB *)
  Database.read db path 0_L 4096_l
  >>= function
    (* TODO: matching on error strings is awful *)
  | Result.Error (`Msg "No such file or directory") -> Lwt.return_none
  | Result.Error (`Msg x) ->
    Lwt.fail (Failure ("Database read failed: "^x))
  | Result.Ok bufs ->
    let len = List.fold_left (+) 0 (List.map Cstruct.len bufs) in
    let buf = Bytes.create len in
    ignore (List.fold_left (fun off cs ->
      let len = Cstruct.len cs in
      Cstruct.blit_to_bytes cs 0 buf off len;
      off + len
    ) 0 bufs);
    Lwt.return (Some buf)

let create_key db path mode data =
  let name = List.(hd (rev path)) in
  let parent = List.(rev (tl (rev path))) in
  Database.create db parent name mode
  >>= function
  | Result.Error (`Msg x) ->
    Lwt.return (Result.Error (`Msg ("create failed: "^x)))
  | Result.Ok () ->
    Database.write db path 0_L (Cstruct.of_string data)
    >>= function
    | Result.Error (`Msg x) ->
      Lwt.return (Result.Error (`Msg ("write failed: "^x)))
    | Result.Ok () -> Lwt.return (Result.Ok ())

let parse data =
  (* TODO: technically paths can include \n and \t, handle those edge cases *)
  let lines = Stringext.split (Bytes.to_string data) ~on:'\n' in
  List.fold_left (fun list line ->
    match Stringext.split ~max:2 line ~on:':' with
    | [] -> list
    | [export_path] -> (export_path, export_path)::list
    | export_path::mount_path::_ -> (export_path, mount_path)::list
  ) [] lines

let print = List.fold_left (fun p (export_path, mount_path) ->
  p^(export_path^":"^mount_path)^"\n"
) ""

let make_mounts = List.fold_left (fun last (export_path, mount_path) ->
  last
  >>= fun list ->
  make ~export_path ~mount_path
  >>= fun mount ->
  Lwt.return (mount::list)
) (Lwt.return [])

let rec connect socket = function
  | 0 -> Lwt.fail (Failure "failed to connect to database")
  | k -> Lwt.catch (fun () ->
    Database.connect "unix" socket ()
  ) (function
    | Unix.Unix_error (Unix.ECONNREFUSED, "connect", _label) ->
      let k = k - 1 in
      Log.warn (fun f ->
        f "database connection attempt failed. Retrying %d more times..." k
      );
      Lwt_unix.sleep 0.1
      >>= fun () ->
      connect socket k
    | Unix.Unix_error (err, "connect", _label) ->
      let msg =
        Printf.sprintf "Couldn't connect to database at %s: %s"
          socket (Unix.error_message err)
      in
      Lwt.fail (Failure msg)
    | exn -> Lwt.fail exn
  )

let load database_socket =
  connect database_socket 10
  >>= function
  | Result.Error (`Msg x) ->
    Lwt.fail (Failure ("Database connection failed: "^x))
  | Result.Ok database ->
    let attempts = 10 in
    let rec try_to_transact = function
      | 0 ->
        Lwt.fail (Failure ("Could not connect to database after "
                           ^string_of_int attempts^" attempts"))
      | x ->
        Database.mkdir database ["branch"; "master"; "transactions"]
          txn_name txn_perm
        >>= function
        | Result.Error (`Msg "No such file or directory") ->
          Lwt_unix.sleep 0.1
          >>= fun () ->
          try_to_transact (x - 1)
        | Result.Error (`Msg x) ->
          Lwt.fail (Failure ("Database transaction start failed: "^x))
        | Result.Ok () ->
          let ctl = txn@["ctl"] in
          Lwt.catch (fun () ->
            let config_path = txn@("rw"::config_path) in
            read_key database config_path
            >>= function
            | Some data ->
              let exports = parse data in
              Lwt_list.filter_p (fun (export,_mount) ->
                Lwt_unix.file_exists export
              ) exports
              >>= fun existing ->
              (if existing <> exports
               then Log.warn (fun f ->
                 let missing =
                   List.filter (fun x -> not (List.mem x existing)) exports
                 in
                 f "osxfs export paths are missing: %s"
                   (String.concat ", " (List.map fst missing))
               ));
              let existing = List.filter (function
                | ("/", "/Mac") -> false
                | _ -> true
              ) existing in
              let existing =
                if List.mem host_docker_app existing
                then existing
                else host_docker_app :: existing
              in
              Lwt.return existing
            | None ->
              let mount_config = print default_mounts in
              create_key database config_path mounts_perm mount_config
              >>= function
              | Result.Error (`Msg x) ->
                Lwt.fail (Failure ("Database set default failed: "^x))
              | Result.Ok () ->
                Database.write database ctl 0_L (Cstruct.of_string "commit")
                >>= function
                | Result.Error (`Msg x) ->
                  Lwt.fail (Failure ("Database transaction commit failed: "^x))
                | Result.Ok () -> Lwt.return default_mounts
          ) (fun exn ->
            Database.write database ctl 0_L (Cstruct.of_string "close")
            >>= function
            | Result.Error (`Msg x) ->
              let err = Printexc.to_string exn in
              Lwt.fail (Failure ("Database transaction abort failed: "
                                 ^err^"\n"^x))
            | Result.Ok () -> Lwt.fail exn
          )
          >>= make_mounts
    in
    try_to_transact attempts
