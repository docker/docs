open Misc

let timestamp =
  let ts = Unix.gettimeofday () in
  let tm = Unix.localtime ts in
  (* example: "20120113-182652" *)
  Printf.sprintf "%04d%02d%02d-%02d%02d%02d"
    (1900 + tm.Unix.tm_year)
    (1    + tm.Unix.tm_mon)
    tm.Unix.tm_mday
    tm.Unix.tm_hour
    tm.Unix.tm_min
    tm.Unix.tm_sec

let dir =
  let dir = "/tmp" / User.unique_id () / timestamp in
  Cmd.mkdir dir;
  at_exit (fun () -> Cmd.exec "rm -rf %s" dir);
  dir

let add_file ?subdir src =
  if Sys.file_exists src then begin
    let dir = match subdir with
      | None -> dir
      | Some p -> dir / p in
    Cmd.mkdir dir;
    Cmd.exec "cp -r \"%s\" \"%s/\"" src dir
  end

let add_string ?subdir key value =
  let logdir = match subdir with None -> dir | Some x -> dir / x in
  let oc = open_out (logdir / key) in
  output_string oc value;
  close_out oc

let collect () =
  add_file "/tmp/com.docker.driver.amd64-linux.log";
  if Sys.file_exists "/var/log/system.log" then (
    let file = dir / "docker-system.log" in
    let filter_expr = String.concat " -o " [
      "-k Sender Seq docker";
      "-k Sender Seq Docker";
      "-k Message Seq docker";
      "-k Message Seq Docker";
    ] in
    let format =
      "$Time $Host $(Sender)[$(Facility)][$(PID)] <$((Level)(str))>: $Message"
    in
    Cmd.exec "syslog -F '%s' %s > \"%s\"" format filter_expr file;
  );
  let () =
    let plist = dir / "group.com.docker.plist" in
    try Cmd.exec "plutil -p %S > %s" (prefs / "group.com.docker.plist") plist
    with e -> Cmd.exec "echo 'ERROR: %S' > %s" (Printexc.to_string e) plist
  in
  if Sys.file_exists "/Library/Preferences/com.apple.alf.plist" then (
    Cmd.exec "defaults read /Library/Preferences/com.apple.alf.plist > %s" (dir / "fw-config");
  );
  Cmd.exec "echo %s > \"%s\"" Version.git (dir / "version");
  Cmd.exec "ps ax > \"%s\"" (dir / "ps-ax.log");
  Cmd.exec "sw_vers > %s" (dir / "sw_vers");
  Cmd.exec "sysctl -a > %s" (dir / "sysctl-a");
  Cmd.exec "date > %s" (dir / "date");
  Cmd.exec "date -u > %s" (dir / "date-u")

let tar outfile =
  let filename = match outfile with
  | None -> dir ^ ".tar.gz"
  | Some x -> x in
  Cmd.exec "tar -C /tmp -cz \"%s\" > \"%s\"" (User.unique_id () / timestamp) filename;
  show "Docker logs are being collected into %s" filename;
  filename
