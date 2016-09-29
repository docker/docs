module Vsock: sig

  val connect: path:string -> ?rcvbuf:int -> ?sndbuf:int -> port:int32 -> unit -> Lwt_unix.file_descr Lwt.t
  (** [connect ~path ?rcvbuf ?sndbuf ~port ()] returns a file descriptor connected
      to [port] in the hyperkit VM, via the hypervisor proxy listening on [path].
      The default values for the receive and send socket buffers are 32KiB. *)
end
