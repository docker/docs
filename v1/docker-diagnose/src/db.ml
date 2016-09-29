open Misc
open Common

let name = "db"
let plist = home / "Library/LaunchAgents/com.docker.db.plist"
let check () =
  check_no_file name plist;
  check_no_lctl name "com.docker.db";
  match Sys.getenv "HOME" with
  | home ->
    check_socket name (Filename.concat home "Library/Containers/com.docker.docker/Data/s40");
    check_ps name "com.docker.db";
    ok name
  | exception Not_found -> error name "HOME not set"
