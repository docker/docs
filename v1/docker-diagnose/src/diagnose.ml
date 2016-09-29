open Misc
open Common

let try_with (name, f) =
  last_errors := [];
  let return () =
    match !last_errors with
    | [] -> { name; success = true ; messages = [] }
    | _  -> { name; success = false; messages = List.rev !last_errors }
  in
  try
    f ();
    return ()
  with e ->
    (* Unhandled exception means the test itself failed *)
    let err = name ^ " check failed with: " ^ (Printexc.to_string e) in
    add_error err;
    return ()

let supported_install () =
  let plist = prefs / "group.com.docker.plist"  in
  let xml =
    Cmd.read_stdout "plutil -extract appVersionHistory xml1 %S -o -" plist
  in
  let _, xml = Ezxmlm.from_string (String.concat "\n" xml) in
  let remove_data =
    List.filter (fun elt -> match elt with `Data _ -> false | _ -> true)
  in
  let kv k v =
    let get str x = Ezxmlm.(member str [x] |> data_to_string) in
    try Some (get "key" k, get "string" v)
    with Ezxmlm.Tag_not_found _ -> None
  in
  let rec version_of = function
    | [] | [_] -> "1.12.0" (* we didn't put the migration code in 1.12.0 *)
    | k::v::t  ->
      match kv k v with
      | Some ("version", v) -> v
      | _                   -> version_of t
  in
  let is_beta v =
    try "beta" = String.sub v (1 + String.rindex v '-') 4
    with Not_found -> false
  in
  let versions =
    xml
    |> Ezxmlm.member  "plist"
    |> Ezxmlm.member  "array"
    |> Ezxmlm.members "dict"
    |> List.map (fun x -> x |> remove_data |> version_of |> is_beta)
  in
  List.for_all (fun x -> x) versions ||     (* only on Beta *)
  List.for_all (fun x -> not x) versions  (* only on stable *)

let os_version () =
  match Cmd.read_stdout "uname" with
  | ["Darwin"] ->
    let product_version =
      try Some (List.hd @@ Cmd.read_stdout "sw_vers -productVersion")
      with _ -> None
    in
    let build_version =
      try Some (List.hd @@ Cmd.read_stdout "sw_vers -buildVersion")
      with _ -> None
    in
    begin match product_version, build_version with
      | None  , None   -> "OS X: unknown version"
      | Some x, None   -> Printf.sprintf "OS X: version %s" x
      | Some x, Some y -> Printf.sprintf "OS X: version %s (build: %s)" x y
      | None  , Some y -> Printf.sprintf "OS X: build %s" y
    end
  | _ -> "unknown OS"

let bundle_version () =
  let bundle = Filename.dirname Sys.argv.(0) / "../../.." in
  let plist = bundle / "Contents" / "Info.plist" in
  if not (Sys.file_exists plist) then "unknown" else
  let plist_in = open_in plist in
  finally
    (fun () ->
      let _, xml = Ezxmlm.from_channel plist_in in
      let dict = List.filter (fun elt -> match elt with `Data _ -> false | _ -> true)
        @@ Ezxmlm.member "dict" @@ Ezxmlm.member "plist" xml in
      let pairs, _ = List.fold_left (fun (pairs, key) elt ->
        match elt with
        | `El((("", "key"), []), [ `Data key ]) -> pairs, key
        | `El((("", _type), []), [ `Data value ]) -> (key, value) :: pairs, ""
        | _ -> pairs, key
      ) ([], "") dict in
      let cFBundleShortVersionString = "CFBundleShortVersionString" in
      if List.mem_assoc cFBundleShortVersionString pairs
      then List.assoc cFBundleShortVersionString pairs
      else "unknown"
    ) (fun () -> close_in plist_in)

let process verbose json upload no outfile =
  if verbose then debug_flag := true;
  if json then json_flag := true;
  if upload then upload_flag := true;
  if no then no_flag := true;
  let os = os_version () in
  let app =
    let s =
      try if supported_install () then "" else "-unsupported"
      with _ -> ""
    in
    Printf.sprintf "version: %s%s (%s)" (bundle_version ()) s Version.git
  in
  show "%s\nDocker.app: %s" os app;
  show "Local time: %s" (String.concat " " @@ Cmd.read_stdout "date");
  show "UTC:        %s" (String.concat " " @@ Cmd.read_stdout "date -u");
  show "Timestamp:  %s" Logs.timestamp;
  show "Running diagnostic tests:";
  let tests = List.map try_with [
      (* Perform the checks most general first, and allow the most specific
         failing check to decide the error class. *)
      "docker-cli", DockerCli.check;
      "moby", Moby.check;
      "driver.amd64-linux", Driver.Amd64_linux.check;
      "vmnetd", Vmnetd.check;
      "osxfs", Osxfs.check;
      "db", Db.check;
      "slirp", Slirp.check;
      "dns", Dns.check;
      "disk", Disk.check;
      "menubar", Menubar.check;
      "env", Env.check;
      "app", App.check;
      "virtualization VT-X", Virtualization.check_vtx;
      "virtualization kern.hv_support", Virtualization.check_kern_hv_support;
      "logs", Logs.collect;
      "app", App.collect;
      "moby-console", Moby.console;
      "moby-syslog", Moby.syslog;
      "moby", Moby.collect;
      "system", System.collect;
    ] in
  let failure = match get_error_class () with
    | None   -> None
    | Some f ->
      Logs.add_string "error-class" f;
      Logs.add_string "errors" (String.concat "\n" !errors);
      show "Most specific failure is: %s" f;
      Some f
  in
  let id = User.unique_id () in
  Logs.add_string "user-id" id;
  let logs = Logs.tar outfile in
  let json = { os; app; id; logs; failure; tests } in
  if !json_flag then show_json json;
  if not !upload_flag && not !no_flag then begin
    Printf.fprintf stdout "Would you like to upload log files? [Y/n]: ";
    flush stdout;
    begin match String.lowercase (input_line stdin) with
      | "" | "y" | "yes" -> upload_flag := true
      | _ -> ()
    end;
    Printf.fprintf stdout "\n";
  end;

  if !upload_flag then begin
    (* upload to S3 *)
    let bucket = "docker-pinata-support" in
    let ctype = "application/json" in
    let date =
      match Cmd.read_stdout "date '+%%a, %%d %%b %%Y %%T %%Z'" with
      |[hd] -> hd
      |_ -> Printf.printf "unable to get current date\n%!"; exit 1 in
    let path = "incoming/2/" / id / Logs.timestamp / "report.tar" in
    Cmd.exec
      "/usr/bin/curl -L -X PUT -T '%s' -H 'Host: %s' -H 'Date: %s' -H 'x-amz-acl: bucket-owner-full-control' -H 'Content-type: %s' https://%s.s3.amazonaws.com/%s"
      logs bucket date ctype bucket path;

      if !json_flag == false then
        Printf.printf "Your unique id is: %s\n" id;
      if !json_flag == false then
        Printf.printf "Please quote this in all correspondence.\n";
  end

open Cmdliner

let verbose =
  let doc = Arg.info ~doc:"Be more verbose." ["v";"verbose"] in
  Arg.(value & flag & doc)

let upload =
  let doc = Arg.info ~doc:"Upload a report for analysis." [ "u"; "upload" ] in
  Arg.(value & flag & doc)

let no =
  let doc = Arg.info ~doc:"Always answer 'No' to interactive questions."
      [ "n"; "no" ] in
  Arg.(value & flag & doc)

let json =
  let doc = Arg.info ~doc:"Output a JSON summary." [ "json" ] in
  Arg.(value & flag & doc)

let outfile =
  let doc = Arg.info ~doc:"Output tarball filename." [ "o"; "outfile" ] in
  Arg.(value & opt (some string) None & doc)

let term =
  let doc = "Docker diagnostic tool." in
  let man = [
    `S "DESCRIPTION";
    `P "$(i, docker-diagnose) is a small tool to check for classic \
        errors when running Docker on OSX, and upload logs for analysis"
  ] in
  Term.(pure process $ verbose $ json $ upload $ no $ outfile),
  Term.info "diagnose" ~version:Version.git ~doc ~man

let () = match Term.eval term with
  | `Error _ -> exit 1
  | _        -> if !errors = [] then () else exit 1
