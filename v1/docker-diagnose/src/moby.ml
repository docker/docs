open Misc
open Common

let driver = "com.docker.driver.amd64-linux"
let console_ring = app / driver / "console-ring"
let vsock_syslog = app / driver / "syslog"
let log_dir = app / driver / "log"
let connect_sock = home / "Library/Containers/com.docker.docker/Data/@connect"

let check () =
  check_file "Moby booted" console_ring;
  ok "Moby booted"

let console () =
  try
    Logs.add_file console_ring
  with e ->
    debug "Unable to find moby console: %s" (Printexc.to_string e)

let syslog () =
  try
    Logs.add_file vsock_syslog
  with e ->
    debug "Unable to find moby syslog: %s" (Printexc.to_string e)

let collect_var_log () =
  List.iter
    (fun log ->
      try
        let path = log_dir / log in
        Logs.add_file ~subdir:"moby/var/log" path
      with e ->
        debug "Unable to upload moby log %s: %s" log (Printexc.to_string e)
    ) (Array.to_list (Sys.readdir log_dir))

let collect_diagnostics () =
  if Sys.file_exists connect_sock then begin
    Cmd.exec ~timeout:30. "echo \"00000003.0000f3a6\" | nc -U %s > %s/diagnostics.tar" connect_sock Logs.dir
  end

let collect () =
  collect_var_log ();
  collect_diagnostics ()
