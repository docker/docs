open Common
open Misc

let name = "slirp"
let home = match Sys.getenv "HOME" with
  | home -> Some home
  | exception Not_found ->
    error name "HOME not set";
    None

let ( / ) = Filename.concat

let send_sigusr1 home =
  let container_dir = home / "Library/Containers/com.docker.docker" in
  let slirp_task = container_dir / "Data/tasks/com.docker.slirp" in
  try
    let ic = open_in slirp_task in
    finally
      (fun () ->
        let json = Ezjsonm.from_channel ic in
        let pid = Ezjsonm.(get_int (find json ["Pid"])) in
        Unix.kill pid Sys.sigusr1
      ) (fun () -> close_in ic)
  with e ->
    debug "Unable to read %s: %s" slirp_task (Printexc.to_string e)

let check () = match home with
  | Some home ->
    check_socket name (home / "Library/Containers/com.docker.docker/Data/s51");
    check_ps name "com.docker.slirp";
    send_sigusr1 home;
    ok name
  | None -> ()
