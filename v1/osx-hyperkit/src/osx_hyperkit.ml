let src =
  let src = Logs.Src.create "osx-hyperkit" ~doc:"Hyperkit interface" in
  Logs.Src.set_level src (Some Logs.Debug);
  src

module Log = (val Logs.src_log src : Logs.LOG)

module Vsock = struct

  open Lwt.Infix

  let connect ~path ?(rcvbuf = 1 lsl 15) ?(sndbuf = 1 lsl 15) ~port () =
    let sock = Lwt_unix.(socket PF_UNIX SOCK_STREAM 0) in
    Lwt_unix.(setsockopt_int sock SO_SNDBUF) sndbuf;
    Lwt_unix.(setsockopt_int sock SO_RCVBUF) rcvbuf;
    Lwt.catch (fun () ->
      Lwt_unix.(connect sock (ADDR_UNIX path))
      >>= fun () ->
      (* There can be only one VM address so it's hardcoded *)
      let address = Cstruct.of_string (Printf.sprintf "00000003.%08lx\n" port) in
      Lwt.catch (fun () ->
        (* will fail with End_of_file if write is short *)
        Lwt_cstruct.complete (Lwt_cstruct.write sock) address
        >>= fun () ->
        Lwt.return sock
      ) (function
        | exn ->
          Log.err (fun f ->
            f "vsock connect write: %s" (Printexc.to_string exn));
          Lwt.fail exn
      )
    ) (function
      | exn ->
        Lwt_unix.close sock
        >>= fun () ->
        Lwt.fail exn
    )
end
