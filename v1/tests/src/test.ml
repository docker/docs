(* config *)

let (/) = Filename.concat

let app_version =
  try Some (Sys.getenv "PINATA_APP_VERSION")
  with Not_found -> None

let app_location =
  try Some (Sys.getenv "PINATA_APP_PATH")
  with Not_found -> None

let cache = "/tmp/docker"

let home = try Sys.getenv "HOME" with Not_found -> "/Users/root"

(* let backend_only = try ignore(Sys.getenv "CI"); true with Not_found -> false *)
let backend_only = false


let containers = home / "Library/Containers/com.docker.docker/Data"
let group_containers = home / "/Library/Group Containers/group.com.docker"

(* misc *)

let date () =
  let now = Unix.gettimeofday () in
  (* Feb 26 07:07:50 *)
  let local = Unix.localtime now in
  let month = match local.Unix.tm_mon with
    | 0 -> "Jan" | 1 -> "Feb" | 2 -> "Mar" | 3 -> "Apr" | 4 -> "May" | 5 -> "Jun"
    | 6 -> "Jul" | 7 -> "Aug" | 8 -> "Sep" | 9 -> "Oct" | 10 -> "Nov" | 11 -> "Dec"
    | _ -> assert false in
  Printf.sprintf "%s %d %02d:%02d:%02d" month local.Unix.tm_mday
    local.Unix.tm_hour local.Unix.tm_min local.Unix.tm_sec

let err fmt = Printf.ksprintf (fun e -> Alcotest.fail (date () ^ " " ^ e)) fmt
let fail cmd fmt = Printf.ksprintf (fun e -> err "%s %s: %s" (date ()) cmd e) fmt
let failv cmd fmt =
  Printf.ksprintf (fun e -> `Error (Printf.sprintf "%s %s: %s" (date ()) cmd e)) fmt

let show fmt =
  Fmt.kstrf (fun str ->
      Fmt.(pf stdout) "%s %a\n%!" (date ()) Fmt.(styled `Blue string) str
    ) fmt

let mutex = Mutex.create ()

let with_lock f =
  Mutex.lock mutex;
  try let r = f () in Mutex.unlock mutex; r
  with e -> Mutex.unlock mutex; raise e

let debug fmt =
  with_lock (fun () ->
      Fmt.kstrf (fun cmd ->
          Fmt.(pf stdout) "%s %a %s\n%!" (date ()) Fmt.(styled `Yellow string) "=>" cmd
        ) fmt
    )

let finally f g =
  try
    let r = f () in
    g ();
    r
  with
  | e ->
    g ();
    raise e

module Thread = struct

  include Thread

  let pick f1 f2 =
    let ch = Event.new_channel () in
    let return x = Event.sync (Event.send ch x) in
    let t1 = Thread.create (fun () -> return (f1 ())) () in
    let t2 = Thread.create (fun () -> return (f2 ())) () in
    let kill () =
      let try_kill t =  try Thread.kill t; Thread.join t with _ -> () in
      try_kill t1;
      try_kill t2;
    in
    let res = Event.sync (Event.receive ch) in
    kill ();
    res

  let with_timeout seconds f =
    let timeout () =
      Thread.delay seconds;
      `Timeout
    in
    let tf () = `Ok (f ()) in
    match pick timeout tf with
    | `Ok x -> x
    | `Timeout -> err "timeout after %.2fs" seconds

end

module Cmd = struct

  let check_exit_status cmd (out, err) status =
    match status with
    | Unix.WEXITED 0   ->
      debug "command successful";
      `Ok (out, err)
    | Unix.WEXITED i   ->
      debug "command failed with code %d" i;
      failv cmd "exit %d" i
    | Unix.WSIGNALED i ->
      debug "command caught signal %d" i;
      if i = Sys.sigkill then fail cmd "timeout" else fail cmd "signal %d" i
    | Unix.WSTOPPED i  ->
      debug "command stopped %d" i;
      fail cmd "stopped %d" i

  let read_lines ?(prefix="") oc =
    let rec aux acc =
      let line =
        try
          let line = input_line oc in
          debug "%s %s" prefix line;
          Some line
        with End_of_file -> None
      in
      match line with
      | Some l -> aux (l :: acc)
      | None   -> List.rev acc
    in
    aux []

  let syscall cmd =
    let env = Unix.environment () in
    (* If 2 threads call Unix.open_process_full then one can
       allocate the pipes for stdout/stderr and then the other
       thread can fork() and capture it. *)
    let oc, ic, ec = with_lock (fun () -> Unix.open_process_full cmd env) in
    let out = ref [] in
    let err = ref [] in
    let stdout_t = Thread.create (fun () ->
      out := read_lines ~prefix:"stdout:" oc
    ) () in
    let stderr_t = Thread.create (fun () ->
      err := read_lines ~prefix:"stderr:" ec
    ) () in
    Thread.join stdout_t;
    Thread.join stderr_t;
    let exit_status = Unix.close_process_full (oc, ic, ec) in
    check_exit_status cmd (!out, !err) exit_status

  let read_outputs ?timeout fmt =
    let f () = Printf.ksprintf syscall fmt in
    match timeout with
    | None   -> f ()
    | Some t -> Thread.with_timeout t f

  let exec ?timeout fmt =
    Printf.ksprintf (fun cmd ->
        debug "%s" cmd;
        match read_outputs ?timeout "%s" cmd with
        | `Ok _    -> ()
        | `Error e -> Alcotest.fail e
      ) fmt

  let read_stdout ?timeout fmt =
    Printf.ksprintf (fun cmd ->
        debug "%s" cmd;
        match read_outputs ?timeout "%s" cmd with
        | `Ok (out, _) -> out
        | `Error e     -> Alcotest.fail e
      ) fmt

  let try_read_stdout ~timeout fmt =
    let t0 = Unix.gettimeofday () in
    Printf.ksprintf (fun cmd ->
        debug "%s (with timeout = %f)" cmd timeout;
        Thread.with_timeout timeout
          (fun () ->
             let rec loop ?error time =
               if time > timeout then match error with
                 | None   -> assert false
                 | Some e -> Alcotest.fail e
               else
               match read_outputs ~timeout:(timeout -. time) "%s" cmd with
               | `Ok (out, _) -> out
               | `Error error ->
                 debug "got: %s, sleeping for 1s (%.2fs/%.2fs)" error time timeout;
                 Thread.delay 1.;
                 loop ~error (Unix.gettimeofday () -. t0)
             in
             loop 0.)
      ) fmt

  let exists s = Sys.command ("which " ^ s) = 0

  let in_dir dir f =
    let pwd = Sys.getcwd () in
    let reset () = if pwd <> dir then Sys.chdir pwd in
    if pwd <> dir then Sys.chdir dir;
    try let r = f () in reset (); r
    with e -> reset (); raise e

  let rmdir path = exec "rm -rf %s" path
  let mkdir path = exec "mkdir -p %s" path

  let ps () = read_stdout "ps ax | grep -i docker"

  let is_running_in ps_lines app =
    let re = Re.compile (Re_glob.globx app) in
    let exists =
      List.exists (fun line -> Re.execp re line) ps_lines
    in
    if exists then debug "%s is running" app
    else debug "%s is not running" app;
    exists

  let is_running app = is_running_in (ps ()) app

  let wait_for_file ?(n=10) file =
    let rec loop remaining = match Sys.file_exists file, remaining with
      | true, _ ->
        debug "%s exists!" file
      | false, 0 ->
        err "%s does not exist after %d seconds" file n
      | false, n ->
        debug "%s does not exist yet, sleeping for 1s" file;
        Thread.delay 1.;
        loop (n-1) in
    loop n

  let wait_for_process ?(n=10) app =
    let rec loop = function
      | 0 -> ()
      | n ->
        if is_running app then
          debug "%s is running!" app
        else (
          debug "%s is not running, sleeping for 1s" app;
          Thread.delay 1.;
          loop (n-1)
        ) in
    loop n

  let touch file = exec "touch %s" file
  let remove file =exec "rm -f %s" file
  let unzip file =  exec "unzip %s" file

  let docker ?ip ?socket ?timeout ~path fmt =
    let h = match ip, socket with
      | None, None -> ""
      | Some ip, _ -> Printf.sprintf "-H tcp://%s:2375" ip
      | None, Some path -> Printf.sprintf "-H unix://%s" path in
    Printf.ksprintf (fun cmd ->
        match timeout with
        | None         ->
          read_stdout "%s/Contents/Resources/bin/docker %s %s" path h cmd
        | Some timeout ->
          try_read_stdout ~timeout "%s/Contents/Resources/bin/docker %s %s"
            path h cmd
      ) fmt

  let docker_run ?timeout ?(volumes=[]) ~path fmt =
    Printf.ksprintf (fun cmd ->
        let volumes =
          List.map (fun (k, v) -> Printf.sprintf "-v %S:%S " k v) volumes
        in
        let volumes = String.concat " " volumes in
        docker ?timeout ~path "run --rm %s%s" volumes cmd
      ) fmt

  let docker_run_sh ?timeout ?(volumes=[]) ~path where lines =
    (* Note the lack of escaping here *)
    let lines =  String.concat ";" lines in
    docker_run ?timeout ~volumes ~path "%s sh -c '%s'" where lines

  let write_file ~output file =
    let oc = open_out_bin file in
    output_string oc output;
    close_out oc

end

(* meat *)

module HockeyApp = struct

  let url = "https://rink.hockeyapp.net/api/2/apps"
  let app_id = "e89aea10ddc530011df5c23127e4fcc1"
  let token = "c59cb6130b24406cb6de08036e67fd1b"

  let curl output fmt =
    Printf.ksprintf (fun path ->
        Cmd.exec "curl -o %s -LH \"X-HockeyAppToken: %s\" %s/%s/app_versions%s"
          output token url app_id path
      ) fmt

  let last_version =
    let version = lazy (
      let versions = cache / "versions" in
      Cmd.mkdir cache;
      curl versions "";
      let all =
        let open Ezjsonm in
        let ic = open_in versions in
        let json = Ezjsonm.from_channel ic in
        close_in ic;
        let version_of_string x =
          let int_of_string x = try int_of_string x with _ -> -1 in
          List.map int_of_string @@ Stringext.split ~on:'.' x in
        List.map (fun dict ->
          string_of_int @@ int_of_float @@ get_float @@ List.assoc "id" dict,
          version_of_string @@ get_string @@ List.assoc "version" dict
        ) (get_list get_dict @@ find json [ "app_versions" ]) in
      (* Default compare on the tuple should order the versions correctly *)
      match List.sort (fun (_, v1) (_, v2) -> compare v2 v1) all with
      | (id, v) :: _ ->
        id
      | []  -> err "unkown application version"
    ) in
    fun () ->
      Lazy.force version

  let download version = curl "Docker.app.zip" "/%s?format=zip" version

end

module Bundle = struct

  let download version =
    let file = cache / version / "Docker.app" in
    if Sys.file_exists file then (
      show "+ Use the cached file.";
    ) else (
      show "+ Download v%s" version;
      Cmd.mkdir (cache / version);
      Cmd.in_dir (cache / version) (fun () ->
          HockeyApp.download version;
          Cmd.unzip "Docker.app.zip";
        );
    );
    file

  let path =
    let path = lazy (
      match app_location with
      | Some file ->
        if not (Sys.file_exists file) || not (Sys.is_directory file) then
          err "%s is not a valid Application Bundle" file;
        show "+ Use %s" file;
        file
      | None      ->
        let version = match app_version with
          | None   -> HockeyApp.last_version ()
          | Some v -> v
        in
        let file = download version in
        show "+ Download v%s" version;
        file
    ) in
    fun () -> Lazy.force path
end


let wait_for_docker () =

  let port_control = home /  "Library/Containers/com.docker.docker/Data/s51" in

  let driver_dir = containers / "com.docker.driver.amd64-linux" in
  let app_file = Bundle.path () in
  Cmd.wait_for_file ~n:10 "/var/run/docker.sock";

  Cmd.wait_for_file ~n:30 port_control;
  Cmd.wait_for_file ~n:30 (driver_dir / "Docker.qcow2");
  Cmd.wait_for_file ~n:30 (driver_dir / "console-ring");

  let (_:string list) = Cmd.docker ~timeout:180. ~path:app_file "ps" in
  ()

module Backend = struct
  (** Used for testing the backend only. Use this in Headless
      environments like under CI. *)

  let com_docker_backend = Bundle.path () / "Contents/MacOS/com.docker.osx.hyperkit.linux"
  let _ =
    if not(Sys.file_exists com_docker_backend) then begin
      Printf.fprintf stderr "Failed to find %s\n%!" com_docker_backend;
      exit 1
    end

  let stdin_r, stdin_w = Unix.pipe ()
  let ic = Unix.out_channel_of_descr stdin_w

  let install () =
    debug "Backend.install";
    output_string ic "install\n";
    flush ic
  let stop () =
    debug "Backend.stop";
    output_string ic "stop\n";
    flush ic
  let start () =
    stop ();
    debug "Backend.start";
    output_string ic "start\n";
    flush ic;
    wait_for_docker ()

end

module GUI = struct

  (* install starts the app and quits right after install steps.
    The app is launched with root privileges to install vmnetd without prompting. *)
  let install ?(quit_after_install=true) () =
    let args =
      "--token=ughetJWhz5aFzf5dgc8qu24Tp --unattended" ^
      (if quit_after_install then " --quit-after-install" else "")
    in
    let app_file = Bundle.path () in
    Cmd.exec "sudo %s/Contents/MacOS/Docker %s" app_file args

  (* install should be called before start to make sure components requiring root 
  privileges are already installed. Because when starting, the apps goes over install steps
  no matter what, to check if everything is ok, and will prompt for root privileges if not. *)
  let start () =
    let args = "--token=ughetJWhz5aFzf5dgc8qu24Tp --unattended" in
    let app_file = Bundle.path () in
    let (_:Thread.t) = Thread.create (fun () ->
      try Cmd.exec "%s/Contents/MacOS/Docker %s" app_file args
      with _ -> ()
    ) () in
    wait_for_docker ()

  let stop () =
    let lines =
      Cmd.read_stdout "ps ax | grep Contents/MacOS/Docke[r] | cut -d' ' -f1"
    in
    let pids = List.map String.trim lines in
    List.iter (fun pid ->
      Cmd.exec "kill -2 %s" pid;
      Cmd.exec "kill -9 %s" pid;
    ) pids

end

module Error = struct

  (* Check if the environment contains some DOCKER_* variables set by
     docker machine *)
  let docker_env () =
    let bad = [
      "DOCKER_HOST";
      "DOCKER_TLS_VERIFY";
      "DOCKER_CERT_PATH";
      "DOCKER_MACHINE_NAME";
    ] in
    List.fold_left
      (fun acc key ->
        try
          let v = Unix.getenv key in
          debug "unexpected environment variable set: %s (= %s)" key v;
          true
        with
        | Not_found -> acc
      ) false bad

end


module App = struct
  let uninstall () =
    Cmd.exec "%s/Contents/MacOS/Docker --uninstall" (Bundle.path ());
    Thread.delay 5.

  let install () =
    if backend_only
    then Backend.install ()
    else GUI.install ()

  let start () =
    if backend_only
    then Backend.start ()
    else GUI.start ()

  let stop () =
    if backend_only
    then Backend.stop ()
    else GUI.stop ()
end

module Pinata = struct
  (* Invoke the pinata command to get/set configuration *)
  let bin = Bundle.path () / "Contents/Resources/bin/pinata.bin"

  type network = [`Hostnet | `Native]

  let get_network () = match String.lowercase @@ String.trim @@ String.concat "\n" @@ Cmd.read_stdout "%s get network" bin with
  | "hostnet" -> `Hostnet
  | "native" | "nat" -> `Native
  | s -> Alcotest.fail ("network is an unrecognized mode: " ^ s)

  let set_network mode =
    App.stop();
    let string = match mode with
      | `Hostnet -> "hostnet"
      | `Native -> "nat" in
    Cmd.exec "%s set network %s" bin string;
    App.start()

end

module Test = struct

  let vmnetd_process = "com.docker.vmnetd"

  let user_process = [
    "Docker.app/Contents/MacOS/com.docker.driver.amd64-linux";
    "com.docker.helper";
  ]

  let check_vmnetd () =
    let ready = Cmd.is_running vmnetd_process in
    if not ready then err "vmnetd is not running, stopping"

  (* partial installs counts *)
  let is_installed () = List.exists (Cmd.is_running_in (Cmd.ps ())) user_process

  let uninstall ?(privileges=false) () =
    show "+ uninstall the app.";
    (*App.stop ();*)
    App.uninstall ();
    (* App.uninstall uninstalls everything already *)
    (* if not privileges then Vmnetd.uninstall (); *)
    if is_installed () then err "app not properly uninstalled"

  let docker_diagnose () =
    let path = Bundle.path () in
    Cmd.exec "PINATA_APP_PATH=%s %s/Contents/Resources/bin/docker-diagnose -n" path path


  let install () =
    show "install the app.";
    App.install ()
    (* docker_diagnose () *)

  let installAndStart () =
    show "+ install and start the app.";
    App.install ();
    App.start ();
    docker_diagnose ()
    (* App.stop () *)


  let docker_ps timeout () =
    show "+ initial docker ps";
    let ps = Cmd.docker ~timeout ~path:(Bundle.path ()) "ps" in
    if ps = [] then Alcotest.fail "empty ps!"

  let iozone ?output_file ?(test_file="f1") ?(max_size="10m") () =
    show "+ iozone";
    let output_dir = match output_file with
      | None   -> Sys.getcwd () / "iotest"
      | Some f -> Filename.dirname f
    in
    let output_file = match output_file with
      | None   -> output_dir / "results.iozone"
      | Some f -> f
    in
    if not (Sys.file_exists output_dir) then Cmd.mkdir output_dir;
    let path = Bundle.path () in
    let output =
      Cmd.docker_run ~path ~volumes:[output_dir, "/iotest"]
        "-i threadx/docker-ubuntu-iozone iozone -a -n 4k -g %s -f %s"
        max_size test_file
      |> String.concat "/"
    in
    Cmd.write_file ~output output_file

  let fsbenchmark () =
    let bench = "v1/perf/fsbenchmark/bench.sh" in
    if Sys.file_exists bench then (
      let dir = Filename.dirname bench in
      let base = Filename.basename bench in
      Cmd.exec "cd %s && ./%s" dir base
    )

  let git_commit () =
    let tmpdir = "/private/tmp/docker.volume" in
    let path = Bundle.path () in
    Cmd.rmdir tmpdir;
    Cmd.mkdir tmpdir;
    let volumes = [ (tmpdir, "/tmp") ] in
    let (_:string list) =
      Cmd.docker_run_sh ~path ~timeout:60. ~volumes "ubuntu" [
        "apt-get install git -y";
        "cd /tmp";
        "mkdir repo";
        "cd repo";
        "git config --global user.name name";
        "git config --global user.email someone@docker.com";
        "git init";
        "echo hello > there";
        "git add there";
        "git commit -m firstcommit"
      ] in
    ()

  let compose_getting_started () =
    let tmpdir = "/private/tmp/docker.compose" in
    Cmd.rmdir tmpdir;
    Cmd.mkdir tmpdir;
    let olddir = Unix.getcwd () in
    finally
      (fun () ->
        Unix.chdir tmpdir;
        let datadir = Filename.concat olddir "v1/tests/cases/compose_getting_started" in
        List.iter
          (fun file ->
            Cmd.exec "cp %s %s" (Filename.concat datadir file) (Filename.concat tmpdir file)
          ) [ "Dockerfile"; "app.py"; "requirements.txt"; "docker-compose.yml" ];
        Cmd.exec "docker build -t web .";
        Cmd.exec "docker-compose stop";
        Cmd.exec "docker-compose up"
      ) (fun () ->
        Unix.chdir olddir
      )

  let compose_hello_world () =
    let tmpdir = "/private/tmp/docker.compose" in
    Cmd.rmdir tmpdir;
    Cmd.mkdir tmpdir;
    let olddir = Unix.getcwd () in
    finally
      (fun () ->
        Unix.chdir tmpdir;
        let datadir = Filename.concat olddir "v1/tests/cases/compose_hello_world" in
        List.iter
          (fun file ->
            Cmd.exec "cp %s %s" (Filename.concat datadir file) (Filename.concat tmpdir file)
          ) [ "docker-compose.yml" ];
        Cmd.exec "docker-compose stop";
        Cmd.exec "docker-compose up"
    ) (fun () ->
      Unix.chdir olddir
    )

  let compose_elasticsearch () =
    let tmpdir = "/private/tmp/docker.compose" in
    Cmd.rmdir tmpdir;
    Cmd.mkdir tmpdir;
    let olddir = Unix.getcwd () in
    finally
      (fun () ->
        Unix.chdir tmpdir;
        let datadir = Filename.concat olddir "v1/tests/cases/compose_elasticsearch" in
        List.iter
          (fun file ->
            Cmd.exec "cp %s %s" (Filename.concat datadir file) (Filename.concat tmpdir file)
          ) [ "docker-compose.yml" ];
        Cmd.exec "docker-compose stop";
        Cmd.exec "docker-compose start"
      ) (fun () ->
        Unix.chdir olddir
      )
end

let run f () =
  try f ()
  with e ->
    Fmt.(pf stderr) "%a %s\n%!"
      Fmt.(styled `Red string) "Error:" (Printexc.to_string e);
    if Error.docker_env () then
      Printf.eprintf
        "Hint: Some DOCKER_* variable are set in the environment. These \
         should be removed to make Docker.app work properly.\n";
    show "Running docker-diagnose:";
    (try Test.docker_diagnose() with _ -> ());
    show "VM console ring:";
    let ring = home / "Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/console-ring" in
    (try Cmd.exec "cat %s | perl -np -e 's/\\0//g'" ring with _ -> Printf.printf "CONSOLE MISSING\n");
    raise e

let minimal_install = [
  "uninstall" , `Quick, run @@ Test.uninstall;
  "install and quit"  , `Quick, run @@ Test.install;
]

let test_gui_install = minimal_install @ [
  "uninstall again", `Quick, run @@ Test.uninstall;
  "start (install and keep running)"  , `Quick, run @@ Test.installAndStart;
  "first docker ps", `Quick, run @@ Test.docker_ps 10.;
  "next docker ps" , `Quick, run @@ Test.docker_ps 1.;
]

let test_install = if backend_only then minimal_install else test_gui_install

let test_benchmarks = [
  "iozone 10m" , `Slow, run @@ Test.iozone ~max_size:"10m";
  "fsbenchmark", `Slow, run @@ Test.fsbenchmark;
]

(* Wrapper which exposes subdirectories containing Makefiles as
   tests. Each Makefile should have a default target which runs
   the test. *)
let shell dir =
  let dir = (Filename.dirname Sys.argv.(0)) / "cases" / dir in
  (* Only consider subdirectories containing Makefiles *)
  let is_test_case name = Sys.file_exists (dir / name / "Makefile") in
  let all = try Sys.readdir dir with _ -> [| |] in
  List.map
    (fun name ->
      name, `Quick, run (fun () ->
        Cmd.exec "make -C %s/%s" dir name
      )
    ) (List.filter is_test_case @@ Array.to_list all)

let test_basic = shell "basic"

let test_volumes = shell "volumes"

let test_network = shell "network"

let compose = shell "compose"

let tests () = [
  "basic"     , test_basic;
  "volumes"   , test_volumes;
  "network"   , test_network;
  "compose"   , compose;
  "benchmarks", test_benchmarks;
]

let () =
  (* Run the install tests to ensure the database state is created *)
  Alcotest.run ~and_exit:false "pinata" [ "install", test_install ];
  (* Now we can query the state to compute the appropriate suite *)
  Alcotest.run ~and_exit:true "pinata" (tests ())
  (*
  Pinata.set_network `Slirp;
  Alcotest.run "pinata" (tests ())
  *)
