open Labels

(* Triage *)

let get_token t =
  match t with
  | None -> Github.Token.of_string ""
  | Some t -> Github.Token.of_string t

let version_number = Re.(seq [
    rep digit;
    str ".";
    rep digit;
    group( seq [
        opt (str ".");
        opt (rep digit);
      ];
      );
  ])

let for_mac_version_re =
  Re.(compile(seq [
      str "Docker for Mac: ";
      opt (str "version: ");
      group ( seq [
          version_number;
          opt (str "-");
          opt (rep alnum);
          opt (str "-");
          opt (rep alnum);
        ]);
    ]))

let get_for_mac_label body =
  let for_mac_matches = Re.exec_opt for_mac_version_re body in
  match for_mac_matches with
  | None -> None
  | Some for_mac_matches -> begin
      match
        try Some (Re.get for_mac_matches 1)
        with _ -> None
      with
      | None -> None;
      | Some s -> ForMac.of_string s
    end

let osx_version_re =
  Re.(compile(seq [
      alt [
        str "macOS: Version ";
        str "OS X: version ";
      ];
      group (
        version_number;
      );
    ]))

let get_osx_label body =
  let osx_matches = Re.exec_opt osx_version_re body in
  match osx_matches with
  | None -> None
  | Some osx_matches -> begin
      match
        try Some (Re.get osx_matches 1)
        with _ -> None
      with
      | None -> None;
      | Some version -> begin
          (* Check if there is a patch version, if yes, strip it *)
          match
            try Some (Re.get osx_matches 2)
            with _ -> None
          with
          | None -> Osx.of_string version
          | Some patch -> begin
              let stop = (String.length version) - (String.length patch) in
              Osx.of_string @@ String.sub version 0 stop
            end
        end
    end

let should_skip l =
  let lst = List.map (fun l ->
      let name = Labels.Status.of_string l.Github_t.label_name in
      match name with
      | None -> false
      | Some _status -> true
    ) l in
  List.fold_left (||) false lst

module IntMap = Map.Make(struct type t = int let compare = compare end)

let triage_issue acc issue =
  let open Github_t in
  (* Check for existing status labels *)
  let labels = issue.issue_labels in
  if should_skip labels
  then Github.Monad.return acc
  else
    (* Check reported version *)
    let for_mac_label = get_for_mac_label issue.issue_body in
    let osx_label = get_osx_label issue.issue_body in
    let computed_labels = {
      for_mac = for_mac_label;
      osx = osx_label;
      (* Initial status is always NeedsTriage *)
      status = Labels.Status.NeedsTriage;
    } in
    Logs.app (fun m ->
      m "Issue %d: %s" issue.issue_number (to_string computed_labels));
    Github.Monad.return (IntMap.add issue.issue_number computed_labels acc)

module Prompt = struct
  type yes_no = [
    | `Yes
    | `No
  ]

  type response = [
    | yes_no
    | `Other of string
  ]

  let yes_no_of_input ~default = function
    | "y" | "yes" -> Some `Yes
    | "n" | "no"  -> Some `No
    | "" -> Some default
    | _other -> None

  let rec yes_no ?(default=`Yes) s =
    let options = match default with
      | `Yes -> " [Y/n]"
      | `No -> " [y/N]"
    in
    Printf.printf "%s%s %!" s options;
    let r = String.(trim (lowercase (input_line stdin))) in
    match yes_no_of_input ~default r with
    | Some yes_no -> yes_no
    | None ->
      Printf.printf "I didn't understand that. Please enter 'y' or 'n'.\n%!";
      yes_no ~default s

end

let triage () token user repo =
  let t = Github.(Monad.(run (
    API.set_token (get_token token)
    >>= fun () ->
    let issues = Issue.for_repo ~user ~repo () in
    Stream.fold triage_issue IntMap.empty issues
    >>= fun labels_to_apply ->
    if IntMap.is_empty labels_to_apply
    then embed Lwt.return_unit
    else
      match Prompt.yes_no "Would you like to apply these labels?" with
      | `No -> embed Lwt.return_unit
      | `Yes ->
        IntMap.fold (fun num l m ->
          m >>= fun () ->
          let labels = Labels.to_string_list l in
          Github.Issue.add_labels ~user ~repo ~num ~labels ()
          >>~ fun _labels -> embed Lwt.return_unit
        ) labels_to_apply (embed Lwt.return_unit)
  ))) in
  Lwt_main.run t

let clean () token user repo =
  let t = Github.(Monad.(run (
    API.set_token (get_token token)
    >>= fun () ->
    let issues = Issue.for_repo ~user ~repo () in
    Stream.fold Nurse_clean.compute_actions [] issues
    >>= List.fold_left Nurse_clean.(fun m action ->
      let open Github_t in
      m >>= fun () ->
      let issue = action.issue in
      Printf.printf "Issue %d: %s\n%s\n"
        issue.issue_number issue.issue_title issue.issue_html_url;
      List.iter (fun change ->
        Printf.printf "%s\n%!" (Nurse_clean.string_of_change change)
      ) action.changes;
      match Prompt.yes_no "Perform these changes?" with
      | `No -> Printf.printf "Not this time\n%!"; return ()
      | `Yes -> clean_issue ~user ~repo action
    ) (embed Lwt.return_unit)
  ))) in
  Lwt_main.run t

let notify () token user repo label filename =
  let t =
    let open Lwt.Infix in
    Lwt_unix.openfile filename [ Lwt_unix.O_RDONLY ] 0
    >>= fun fd ->
    let ic =
      Lwt_io.of_fd ~close:(fun () -> Lwt_unix.close fd) ~mode:Lwt_io.input fd
    in
    Lwt_io.read ic
    >>= fun body ->
    Lwt_io.close ic
    >>= fun () ->
    Printf.printf
      "For issues with label '%s' I will append the following comment:\n%!"
      label;
    Printf.printf "\n%s\n%!" body;

    Github.(Monad.(run (
      let open Github_t in
      API.set_token (get_token token)
      >>= fun () ->
      Stream.iter (fun issue ->
        (* Check for existing status labels *)
        let labels =
          List.map (fun l -> l.Github_t.label_name) issue.issue_labels
        in
        if List.mem label labels then begin
          Printf.printf "Issue %d: %s\n%s\n"
            issue.issue_number issue.issue_title issue.issue_html_url;
          match Prompt.yes_no "Append comment?" with
          | `Yes ->
            Printf.printf "OK, commenting\n%!";
            Issue.create_comment ~user ~repo ~num:issue.issue_number ~body ()
            >>~ fun comment ->
            Printf.printf "%s\n" comment.issue_comment_html_url;
            return ()
          | `No ->
            Printf.printf "Not this time\n%!";
            return ()
        end else return ()
      ) (Issue.for_repo ~state:`All ~user ~repo ())
    ))) in
  Lwt_main.run t

let analyse () diagnostic_id timestamp =
  let open Lwt.Infix in
  let t =
    Cache.get diagnostic_id timestamp
    >>= function
    | Result.Error (`Msg m) ->
      Logs.err (fun f -> f "Failed to download archive from S3: %s" m);
      Lwt.return ()
    | Result.Ok filename ->
      Archive.openarchive filename
      >>= function
      | Result.Ok t ->
        Logs.app (fun f -> f "Analysing report with diagnostic ID %s and timestamp %s"
                     diagnostic_id timestamp
                 );
        Archive.analyse t
        >>= fun symptoms ->
        Lwt_list.iter_s
          (fun symptom ->
             let md = Archive.Symptom.to_markdown symptom in
             Lwt_io.write Lwt_io.stdout md
          ) symptoms
        >>= fun () ->
        Archive.close t
      | Result.Error (`Msg m) ->
        Logs.err (fun f -> f "Failed to open archive: %s" m);
        Lwt.return () in
  Lwt_main.run t

(* Command line interface *)

open Cmdliner

(* Logging *)

let setup_log style_renderer level =
  Fmt_tty.setup_std_outputs ?style_renderer ();
  Logs.set_level level;
  Logs.set_reporter (Logs_fmt.reporter ());
  ()

let setup_log =
  Term.(const setup_log $ Fmt_cli.style_renderer () $ Logs_cli.level ())

(* Commands *)

let github_token =
  let doc = "Github API Token." in
  Arg.(value & opt (some string) None & info ["token"] ~docv:"TOKEN" ~doc)

let github_owner =
  let doc = "Github Repository Owner." in
  Arg.(required & pos 0 (some string) None & info [] ~docv:"OWNER" ~doc)

let github_repo =
  let doc = "Github Repository Name." in
  Arg.(required & pos 1 (some string) None & info [] ~docv:"REPO" ~doc)

let triage_cmd =
  let doc = "Triage bugs" in
  let man = [`S "DESCRIPTION"; `P "Triages bugs";] in
  Term.(const triage $ setup_log $ github_token $ github_owner $ github_repo),
  Term.info "triage" ~doc ~man

let notify_cmd =
  let label =
    let doc = "Issue label" in
    Arg.(required & pos 2 (some string) None & info [] ~docv:"LABEL" ~doc) in
  let filename =
    let doc = "File containing comment to add" in
    Arg.(required & pos 3 (some file) None & info [] ~docv:"FILENAME" ~doc) in
  let doc = "Append a comment to all issues with a specific tag" in
  let man = [
    `S "DESCRIPTION";
    `P "Given a label and a filename, iterate over every issue (open and closed)\
       containing the tag and append the file's contents as a new comment.";
    `P "You will br prompted to confirm every comment.";
  ] in
  Term.(const notify $ setup_log $ github_token $ github_owner $ github_repo $ label $ filename),
  Term.info "notify" ~doc ~man

let analyse_cmd =
  let diagnostic_id =
    let doc = "Diagnostic ID recorded in the tracker" in
    Arg.(value & pos 0 string "F4466F9B-150B-4082-A6F5-6BA8FA0F085B" & info [] ~docv:"ID" ~doc) in
  let timestamp =
    let doc = "Timestamp of the report" in
    Arg.(value & pos 1 string "20160802-140735" & info [] ~docv:"TIMESTAMP" ~doc) in
  let doc = "Look for well-known errors in diagnostic tarballs" in
  let man = [
    `S "DESCRIPTION";
    `P "This is a temporary command which should be merged with triage later"
  ] in
  Term.(const analyse $ setup_log $ diagnostic_id $ timestamp),
  Term.info "analyse" ~doc ~man

let clean_cmd =
  let doc = "Clean bugs" in
  let man = [`S "DESCRIPTION"; `P "Fixes labels and closes stale issues";] in
  Term.(const clean $ setup_log $ github_token $ github_owner $ github_repo),
  Term.info "clean" ~doc ~man

let default_cmd =
  let doc = "assists in the triaging and care of bugs" in
  let man = [ `S "BUGS"; `P "Open an issue on GitHub";] in
  Term.(ret (const (fun _ -> `Help (`Plain, None)) $ const ())),
  Term.info "nurse" ~version:"0.0.1" ~doc ~man

let cmds = [
  triage_cmd;
  analyse_cmd;
  notify_cmd;
  clean_cmd;
]

(* Entrypoint *)

let () = match Term.eval_choice default_cmd cmds with
  | `Error _ -> exit 1
  | _ -> exit (if Logs.err_count () > 0 then 1 else 0)
