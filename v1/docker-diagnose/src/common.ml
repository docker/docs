open Misc

let results_to_upload : (string, (string * string) list) Hashtbl.t = Hashtbl.create 7

(* Most specific one-line error description. Two bugs with the same description
   will be grouped together. *)
let error_class = ref None

let upload_flag = ref false
let no_flag = ref false

let ok_str  = "[OK]     "
let err_str = "[ERROR]  "

let errors : string list ref = ref []
let last_errors : string list ref = ref []

let ok name =
  if not !json_flag && !last_errors = [] then
    Fmt.(pf stdout) "%a %s\n%!" Fmt.(styled `Green string) ok_str name

let add_error err =
  errors := err :: !errors;
  last_errors := err :: !last_errors

let error ?e name fmt =
  Fmt.kstrf (fun str ->
      let str' =
        str ^ (match e with None -> "" | Some e -> ": " ^ Printexc.to_string e)
      in
      if not !json_flag then
        Fmt.(pf stdout) "%a %-16s %s\n%!"
          Fmt.(styled `Red string) err_str name str';
      add_error str;
      error_class := Some str';
    ) fmt

let get_error_class () = !error_class

let check_file name plist =
  try
    if not (Sys.file_exists plist) then error name "%s does not exist" plist;
  with e ->
    error ~e name "%s is missing" plist

let check_no_file name plist =
  try
    if Sys.file_exists plist then error name "%s should not exist!" plist;
  with e ->
    error ~e name "%s is present" plist

let check_is_dir name dir =
  try
    if not(Sys.is_directory dir) then error name "%s exists but is not a directory" dir
  with
  | e ->
    error ~e name "%s does not exist" dir

let is_running_in lines cmd =
  let re = Re.compile (Re_glob.globx cmd) in
  List.exists (fun line -> Re.execp re line) lines

let check_running_in name lines cmd =
  try
    if is_running_in lines cmd then debug "%s is running" cmd
    else error name "%s is not running" cmd
  with e ->
    error ~e name "%s is not running" cmd

let check_not_running_in name lines cmd =
  try
    if not (is_running_in lines cmd) then debug "%s is not running" cmd
    else error name "%s is running" cmd
  with e ->
    error ~e name "%s is running" cmd

let check_socket name socket =
  try
    begin
      check_file name socket;
      let s = Unix.socket Unix.PF_UNIX Unix.SOCK_STREAM 0 in
      finally
        (fun () ->
          try
            Unix.connect s (Unix.ADDR_UNIX(socket))
          with
            | Unix.Unix_error(Unix.EACCES, _, _) ->
              error name "Permission denied (EACCES) connecting to %s: check permissions" socket
            | Unix.Unix_error(Unix.ECONNREFUSED, _, _) ->
              error name "Connection refused (ECONNREFUSED) connecting to %s: check if service is running" socket
            | Unix.Unix_error(m, _, _) ->
              error name "Unexpected error (%s) connecting to %s" (Unix.error_message m) socket
            | e ->
              error name "Unexpected error (%s) connecting to %s" (Printexc.to_string e) socket
        ) (fun () -> Unix.close s)
    end
  with e ->
    error ~e name "connection to %s refused" socket

let ps = lazy (Cmd.read_stdout "ps ax | grep docker -i")
let launchctl_list = lazy (Cmd.read_stdout "/bin/launchctl list | grep docker -i || echo No services registered with launchd")
let check_ps name = check_running_in name (Lazy.force ps)
let check_no_ps name = check_not_running_in name (Lazy.force ps)
let test_ps = is_running_in (Lazy.force ps)
let check_no_lctl name service =
  try
    if not (is_running_in (Lazy.force launchctl_list) service)
    then debug "%s is not registered with launchd" service
    else error name "%s is registered with launchd" service
  with
  | e ->
    error ~e name "%s is still registered with launchd" service

let check_cmd name cmd =
  try
    if not (Cmd.exists cmd) then error name "cannot find %s" cmd
  with e ->
    error ~e name "%s is missing" cmd

type test = {
  name    : string;
  success : bool;
  messages: string list;
}

let json_of_message m =
  Printf.sprintf
    "\n        %S" m

let json_of_test t =
  Printf.sprintf
    "    %S: {\n\
    \      \"success\": %b,\n\
    \      \"message\": [%s]\n\
    \    }"
    t.name t.success
    (String.concat "," (List.map json_of_message t.messages)
     ^ if t.messages = [] then "" else "\n      ")

type json = {
  os     : string;
  app    : string;
  id     : string;
  logs   : string;
  failure: string option;
  tests  : test list;
}

let show_json json =
  let failure = match json.failure with
    | None -> ""
    | Some f -> Printf.sprintf "\"failure\": %S,\n" f
  in
  Fmt.(pf stdout)
    "{\n\
     \"os\": %S,\n\
     \"app\": %S,\n\
     \"id\": %S,\n\
     \"logs\": %S,\n\
     %s\
     \"tests\": {\n\
     %s\n\
    \  }\n\
     }\n%!"
    json.os json.app json.id json.logs failure
    (String.concat ",\n" @@ List.map json_of_test json.tests)
