open Misc
open Common

module Amd64_linux = struct
  let name = "driver.amd64-linux"
  let old_plist = home / "Library/LaunchAgents/com.docker.xhyve.plist"
  let plist = home / "Library/LaunchAgents/com.docker.driver.amd64-linux.plist"
  let driver = app / "com.docker.driver.amd64-linux"
  let container_dir = home / "Library/Containers/com.docker.docker"

  let check () =
    check_no_file name plist;
    check_no_file name old_plist;
    check_no_lctl name "com.docker.driver.amd64-linux";
    check_ps name "com.docker.driver.amd64-linux -db";
    check_is_dir name container_dir;
    ok name
end
