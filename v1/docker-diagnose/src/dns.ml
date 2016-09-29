open Misc

let check () =
  Logs.add_file ~subdir:"etc" "/etc/resolv.conf";
  Cmd.exec ~timeout:30. "scutil --dns > %s/scutil\\ --dns.stdout 2> %s/scutil\\ --dns.stderr"
    Logs.dir Logs.dir
