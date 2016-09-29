open Misc
open Common

let name = "osxfs"
let plist = home / "Library/LaunchAgents/com.docker.osxfs.plist"
let check () =
  check_no_file name plist;
  check_no_lctl name "com.docker.osxfs";
  check_ps name "com.docker.osxfs";
  ok name
