open Misc
open Common

let name = "vmnetd"
let plist = "/Library/LaunchDaemons/com.docker.vmnetd.plist"

let check () =
  check_socket name "/var/tmp/com.docker.vmnetd.socket";
  check_no_file name (home / plist);
  if test_ps "/Library/PrivilegedHelperTools/com.docker.vmnetd" then (
    (* normal code-path *)
    check_file name plist;
    ok name
  ) else (
    (* FIXME: not sure this is still needed *)
    (* test code-path *)
    debug "vmnetd in used in test mode.";
    check_ps name "/private/tmp/vmnetd/com.docker.vmnetd";
  );
