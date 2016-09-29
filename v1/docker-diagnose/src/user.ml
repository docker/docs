open Misc

let full_name () =
  with_temp_file
    (fun name ->
      Cmd.exec "osascript -e \"long user name of (system info)\" > %s" name;
      String.trim (string_of_file name)
    )

(* A unique persistent id so we can associate bugreports together *)
let unique_id () =
  let path = app / "user.id" in
  if not (Sys.file_exists path) then Cmd.exec "uuidgen > \"%s\"" path;
  String.trim (string_of_file path)

let email () = "unknown" (* FIXME *)
