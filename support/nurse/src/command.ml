open Lwt.Infix

let errorf fmt =
  Printf.ksprintf (fun s ->
    Result.Error (`Msg s)
  ) fmt

let run path args =
  let args = Filename.basename path :: args in
  Logs.info (fun f -> f "%s" (String.concat " " args));
  Lwt_process.with_process_in (path, Array.of_list args)
    (fun p ->
      Lwt_io.read p#stdout
      >>= fun output ->
      p#status
      >>= function
      | Unix.WEXITED 0 ->
        Lwt.return (Result.Ok output)
      | Unix.WEXITED n ->
        Lwt.return (Result.Error (`Exit n))
      | Unix.WSIGNALED n | Unix.WSTOPPED n ->
        Lwt.return (errorf "Caught signal %d while running %s" n path)
    )
