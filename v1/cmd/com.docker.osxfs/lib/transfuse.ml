
open Lwt

module Profuse = Profuse_7_23

module Log = Log.Log

let host = Profuse.Host.linux_4_0_5

module type CONF = sig
  val max_write : int
  val connect_path : string
  val address : string
end

let fuse_port = 0x5f5l

let vsock_connect path port =
  Lwt.catch (fun () ->
    Osx_hyperkit.Vsock.connect ~path ~port ()
  ) (function
    | Unix.Unix_error (err, "connect", label) ->
      let msg =
        Printf.sprintf "Couldn't connect to vsock at %s on port %ld: %s"
          path port (Unix.error_message err)
      in
      Lwt.fail (Failure msg)
    | exn -> Lwt.fail exn
  )

let inotify_channel connect_path =
  let open Lwt.Infix in
  vsock_connect connect_path fuse_port
  >>= fun sock ->
  (* TODO: short write/socket conditions *)
  Lwt_unix.write sock (Bytes.of_string "e") 0 1
  >>= fun n ->
  (if n <> 1 then failwith "Couldn't write event channel prelude");
  let buf = Bytes.create 16 in
  Lwt_unix.read sock buf 0 16
  >>= fun n ->
  let actor_pid = Unsigned.UInt32.of_string
      (Bytes.to_string (Bytes.trim (Bytes.sub buf 0 n)))
  in
  Lwt.return (sock, actor_pid)

let listen_forever address handler =
  Osx_socket.listen address
  >>= fun sock ->
  let rec accept () =
    Lwt_unix.accept sock
    >>= fun (client_fd, client_addr) ->
    let client_fd = Lwt_unix.unix_file_descr client_fd in
    let ctl = Control.make client_fd client_addr in
    Lwt.async (Control.service_connection handler ctl);
    accept ()
  in
  accept ()

let connect_to_fuse_sock ~mnt ~argv ~connect_path =
  vsock_connect connect_path fuse_port
  >>= fun sock ->
  let opts = "m" ^ String.concat "\000" (Array.to_list argv @ [mnt]) in
  let opts = Bytes.of_string opts in
  let sz = Bytes.length opts in
  (* TODO: short writes *)
  Lwt_unix.write sock opts 0 sz
  >>= fun n ->
  (if sz <> n then failwith "Couldn't write FUSE mount options");
  Lwt.return (Lwt_unix.unix_file_descr sock)

module Mount(Conf : CONF)(FS:Fuse_lwt.FS_LWT)(IO: Fuse_lwt.IO_LWT)
  : Fuse.MOUNT_IO
    with module IO = IO
     and type t = FS.t
= struct
  module IO = IO
  type t = FS.t

  (* TODO: Should the profuse lwt layer handle partial read/write
     instead of read_exactly/write_exactly above? *)

  let mount negotiate ~argv ~mnt state =
    let open IO in
    connect_to_fuse_sock ~mnt ~argv ~connect_path:Conf.connect_path
    >>= fun fd ->
    let write_mutex = Lwt_mutex.create () in
    let socket = Fuse_lwt.new_socket
        ~read:(fun sz ->
          Socket.read_message fd
          >>= fun carray ->
          let len = Ctypes.CArray.length carray in
          if len > sz
          then
            let msg =
              Printf.sprintf "Read message of size %d but max is %d" len sz
            in
            Lwt.fail (Failure msg)
          else Lwt.return carray
        )
        ~write:(fun p len ->
          Lwt_mutex.lock write_mutex
          >>= fun () ->
          Socket.write_exactly fd p len
          >>= fun () ->
          Lwt_mutex.unlock write_mutex;
          Lwt.return len
        )
        ~nwrite:(fun p len -> Lwt.return 0)
        ~nread:(fun () -> Lwt.return Unsigned.UInt32.zero)
    in
    let id = Fuse_lwt.socket_id socket in
    let max_write = Conf.max_write in
    let init_fs = Profuse.({
      mnt; id; unique=Unsigned.UInt64.zero; version=(0,0); host;
      max_readahead=0; max_write; flags=Profuse.Flags.empty;
    }) in
    IO.In.read init_fs ()
    >>= fun req ->

    Profuse.(In.(match req with
      | { pkt=Message.Init pkt } -> negotiate pkt req state
      | { hdr } ->
        let opcode = Opcode.(
          to_string (of_uint32 (Ctypes.getf hdr Hdr.T.opcode))
        ) in
        let msg =
          Printf.sprintf "Unexpected opcode %s <> FUSE_INIT" opcode
        in
        IO.fail (ProtocolError (init_fs, msg))
    ))

end
