val listen: string -> Lwt_unix.file_descr Lwt.t
(** [listen spec] produces a listening socket by either binding and listening
    or by being passed a socket from the parent program. *)
