open Misc
open Astring

let qcow2 = home / "Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/Docker.qcow2"

let bundle = Filename.dirname Sys.argv.(0) / "../../.."
let qemu_img = bundle / "Contents/MacOS/qemu-img"

let name = "disk"

let check () =
  if Sys.file_exists qcow2 then begin
    (* Since we run with lazy refcounts for performance we expect lots of
       refcounting errors. Ideally qemu-img check would not produce these when
       the image is marked as dirty. *)
    Cmd.mkdir Logs.dir;
    (* Avoid `Cmd.exec` because it will print a scary error on non-zero exit
       code: unfortunately this is expected at least some of the time *)
    let stderr = Logs.dir ^ "/qemu-img-check.stderr" in
    let stdout = Logs.dir ^ "/qemu-img-check.stdout" in
    ignore(Unix.system (Printf.sprintf "%s check %s > %s 2> %s"
      qemu_img qcow2 stdout stderr));
    (* The stderr output contains one line per error. We ignore all errors which
       include the text "refcount=0". *)
    if not (Sys.file_exists stderr)
    then Common.error name "Failed to find qemu-img in bundle"
    else begin
      let ic = open_in stderr in
      finally
        (fun () ->
          let rec count_unexpected_errors total =
            match input_line ic with
            | line ->
              begin match String.find_sub ~sub:"No such file or directory" line,
                    String.find_sub ~sub:"refcount=0" line  with
              | Some _, _
              | None, Some _ -> count_unexpected_errors total
              | _, _ -> Printf.fprintf Pervasives.stderr "line %s\n%!" line;count_unexpected_errors (total + 1)
              end
            | exception End_of_file -> total in
          let unexpected_errors = count_unexpected_errors 0 in
          if unexpected_errors > 0
          then Common.error name "Docker.qcow2 has unexpected errors: either repair offline or reset to factory defaults"
          else Common.ok name
        ) (fun () ->
          close_in ic
        )
    end
  end else begin
    Common.error name "Docker.qcow2 missing: the VM has never been started"
  end
