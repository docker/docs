let () = Fmt_tty.setup_std_outputs ()
let debug_flag = ref false
let json_flag = ref false
let (/) = Filename.concat

let finally f g =
  try
    let r = f () in
    g ();
    r
  with
  | e ->
    g ();
    raise e

let fail fmt =
  Fmt.kstrf (fun e ->
      if not !json_flag then
        Fmt.(pf stderr) "%a %s\n%!" Fmt.(styled `Red string) "Error" e;
      failwith e
    ) fmt

let failv fmt = Printf.ksprintf (fun e -> `Error e) fmt

let debug fmt =
  Fmt.kstrf (fun cmd ->
      if not !json_flag && !debug_flag then
        Fmt.(pf stdout) "%a %s\n%!" Fmt.(styled `Yellow string) "=>" cmd
    ) fmt

let show fmt =
  Fmt.kstrf (fun str ->
      if not !json_flag then
        Fmt.(pf stdout) "%a\n%!" Fmt.(styled `Bold string) str
    ) fmt

(* misc *)

module Thread = struct

  include Thread

  let pick f1 f2 =
    let ch = Event.new_channel () in
    let protect f =
      try f () with e -> Thread.exit (); failwith (Printexc.to_string e)
    in
    let return x = Event.sync (Event.send ch @@ protect x) in
    let t1 = Thread.create (fun () -> return f1) () in
    let t2 = Thread.create (fun () -> return f2) () in
    let kill () =
      let try_kill t =  try Thread.kill t; Thread.join t with _ -> () in
      try_kill t1;
      try_kill t2;
    in
    let res = Event.sync (Event.receive ch) in
    kill ();
    res

  let with_timeout seconds cmd f =
    let timeout () =
      Thread.delay seconds;
      `Timeout
    in
    let tf () = `Ok (f ()) in
    match pick timeout tf with
    | `Ok x -> x
    | `Timeout -> fail "%s: timeout after %.2fs" cmd seconds

end

module Cmd = struct

  let check_exit_status cmd (out, err) status =
    if out <> [] then debug "stdout:\n%s\n" (String.concat "\n" out);
    if err <> [] then debug "stderr:\n%s\n" (String.concat "\n" out);
    match status with
    | Unix.WEXITED 0   -> `Ok (out, err)
    | Unix.WEXITED i   -> failv "%s: exit %d" cmd i
    | Unix.WSIGNALED i ->
      if i = Sys.sigkill then fail "%s: timeout" cmd
      else fail "%s: signal %d" cmd i
    | Unix.WSTOPPED i  -> fail "%s: stopped %d" cmd i

  let read_lines oc =
    let rec aux acc =
      let line =
        try Some (input_line oc)
        with End_of_file -> None
      in
      match line with
      | Some l -> aux (l :: acc)
      | None   -> List.rev acc
    in
    aux []

  let syscall cmd =
    let env = Unix.environment () in
    let oc, ic, ec = Unix.open_process_full cmd env in
    let out = read_lines oc in
    let err = read_lines ec in
    debug "%s" cmd;
    let exit_status = Unix.close_process_full (oc, ic, ec) in
    check_exit_status cmd (out, err) exit_status

  let read_outputs ?timeout cmd =
    let f () = syscall cmd in
    match timeout with
    | None   -> f ()
    | Some t -> Thread.with_timeout t cmd f

  let exec ?timeout fmt =
    Printf.ksprintf (fun cmd ->
        match read_outputs ?timeout cmd with
        | `Ok _    -> ()
        | `Error e -> fail "exec: %s" e
      ) fmt

  let read_stdout ?timeout fmt =
    Printf.ksprintf (fun cmd ->
        match read_outputs ?timeout cmd with
        | `Ok (out, _) -> out
        | `Error e     -> fail "exec: %s" e
      ) fmt

  let try_read_stdout ~timeout fmt =
    let t0 = Unix.gettimeofday () in
    Fmt.kstrf (fun cmd ->
        Thread.with_timeout timeout cmd
          (fun () ->
             let rec loop ?error time =
               if time > timeout then match error with
                 | None   -> assert false
                 | Some e -> fail "%s: got timeout (%s)" cmd e
               else
                 match read_outputs ~timeout:(timeout -. time) cmd with
                 | `Ok (out, _) -> out
                 | `Error error ->
                   debug "got: %s, sleeping for 1s (%.2fs/%.2fs)" error time timeout;
                   Thread.delay 1.;
                   loop ~error (Unix.gettimeofday () -. t0)
             in
             loop 0.)
      ) fmt

  let exists cmd =
    let cmd = Printf.sprintf "sh -c \"command -v '%s'\" > /dev/null 1>&2" cmd in
    match syscall cmd with
    | `Ok _    -> true
    | `Error _ -> false

  let in_dir dir f =
    let pwd = Sys.getcwd () in
    let reset () = if pwd <> dir then Sys.chdir pwd in
    if pwd <> dir then Sys.chdir dir;
    try let r = f () in reset (); r
    with e -> reset (); raise e

  let mkdir path = exec "mkdir -p %s" path

  let wait_for_file ?(n=10) file =
    let rec loop = function
      | 0 -> ()
      | n ->
        if Sys.file_exists file then
          debug "%s exists!" file
        else (
          debug "%s does not exists, sleeping for 1s" file;
          Thread.delay 1.;
          loop (n-1)
        ) in
    loop n

  let touch file = exec "touch %s" file
  let remove file =exec "rm -f %s" file

  (* wait_for is a file which should appear *)
  let open_file ?wait_for file =
    let dot_test = file / "Contents" / "Resources" / ".test" in
    if Sys.file_exists dot_test then debug ".test exists"
    else debug ".test does not exist";
    exec "open %s" file;
    Thread.join (Thread.create (fun () ->
        match wait_for with
        | None      -> ()
        | Some file -> wait_for_file ~n:10 file
      ) ())

  let unzip file =  exec "unzip %s" file

  let docker ?(timeout=10.) ?path fmt  =
    let path = match path with
      | None   -> ""
      | Some p -> p / "Contents/Resources/bin/"
    in
    Printf.ksprintf (fun cmd ->
        try_read_stdout ~timeout "%sdocker %s" path cmd
      ) fmt

  let write_file ~output file =
    let oc = open_out_bin file in
    output_string oc output;
    close_out oc

end

let startswith prefix x =
  let prefix' = String.length prefix and x' = String.length x in
  prefix' < x' && (String.sub x 0 prefix' = prefix)

let app_location =
  try Some (Sys.getenv "PINATA_APP_PATH")
  with Not_found -> None

let home = try Sys.getenv "HOME" with Not_found -> fail "$HOME is not set"

let app = home / "Library/Containers/com.docker.docker/Data"

let prefs =
  home / "/Library/Group Containers/group.com.docker/Library/Preferences"

let with_temp_file f =
  let name = Filename.temp_file "docker-diagnose" "tmp" in
  let r = f name in
  Unix.unlink name;
  r

let string_of_file path =
  let fd = Unix.openfile path [ Unix.O_RDONLY ] 0 in
  let buffer = Buffer.create 16 in
  let bytes = Bytes.create 1024 in
  let rec loop () =
    let n = Unix.read fd bytes 0 (Bytes.length bytes) in
    if n = 0 then () else begin
      (* Inefficiently remove NULL bytes *)
      let sub = Bytes.sub bytes 0 n in
      let sub = String.concat "" (Stringext.split sub ~on:'\000') in
      Buffer.add_bytes buffer sub;
      loop ()
    end in
  loop ();
  Unix.close fd;
  let buf = Buffer.contents buffer in
  let max_len = 50000 in
  let buf_len = String.length buf in
  if buf_len > max_len then (
    debug "Truncating %s to %d bytes" path max_len;
    String.sub buf (buf_len - max_len) max_len) else buf
