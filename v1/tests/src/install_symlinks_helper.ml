
module Init = struct

  cstruct msg {
      uint8_t magic[5];   (* VMN3T *)
      uint32_t version;   (* 1 *)
      uint8_t commit[40];
    } as little_endian

  type t = {
    magic: string;
    version: int32;
    commit: string;
  }

  let sizeof = sizeof_msg

  let default = {
    magic = "VMN3T";
    version = 1l;
    commit = "0123456789012345678901234567890123456789";
  }

  let marshal t rest =
    set_msg_magic t.magic 0 rest;
    set_msg_version rest t.version;
    set_msg_commit t.commit 0 rest;
    Cstruct.shift rest sizeof_msg

  let unmarshal rest =
    let magic = Cstruct.to_string @@ get_msg_magic rest in
    let version = get_msg_version rest in
    let commit = Cstruct.to_string @@ get_msg_commit rest in
    let rest = Cstruct.shift rest sizeof_msg in
    `Ok ({ magic; version; commit }, rest)
end

module Command = struct

  cstruct msg {
      uint8_t command;
    } as little_endian

  type t =
    | Ethernet of string (* 36 bytes *)
    | Uninstall
    | Install_symlinks
    | Bind_ipv4 of Ipaddr.V4.t * int * bool

  let sizeof = sizeof_msg + 36

  let marshal t rest = match t with
    | Ethernet uuid ->
      set_msg_command rest 1;
      let rest = Cstruct.shift rest sizeof_msg in
      Cstruct.blit_from_string uuid 0 rest 0 (String.length uuid);
      Cstruct.shift rest (String.length uuid)
    | Uninstall ->
      set_msg_command rest 2;
      Cstruct.shift rest sizeof_msg
    | Install_symlinks ->
      set_msg_command rest 3;
      Cstruct.shift rest sizeof_msg
    | Bind_ipv4 (ip, port, stream) ->
      set_msg_command rest 6;
      let rest = Cstruct.shift rest sizeof_msg in
      Cstruct.LE.set_uint32 rest 0 (Ipaddr.V4.to_int32 ip);
      let rest = Cstruct.shift rest 4 in
      Cstruct.LE.set_uint16 rest 0 port;
      let rest = Cstruct.shift rest 2 in
      Cstruct.set_uint8 rest 0 (if stream then 0 else 1);
      Cstruct.shift rest 1

  let unmarshal rest =
    match get_msg_command rest with
    | 1 ->
      let uuid = Cstruct.(to_string (sub rest 1 36)) in
      let rest = Cstruct.shift rest 37 in
      `Ok (Ethernet uuid, rest)
    | 2 ->
      `Ok (Uninstall, Cstruct.shift rest 1)
    | 3 ->
      `Ok (Install_symlinks, Cstruct.shift rest 1)
    | n -> `Error (`Msg (Printf.sprintf "Unknown command: %d" n))
end

module Vmnet_client = struct
  let error_of_failure f = Lwt.catch f (fun e -> Lwt.return (`Error (`Msg (Printexc.to_string e))))
  open Lwt.Infix

  type t = {
    fd: Lwt_unix.file_descr;
  }

  module Infix = struct
    let ( >>= ) m f = m >>= function
      | `Ok x -> f x
      | `Error x -> Lwt.return (`Error x)
  end

  let of_fd fd =
    let buf = Cstruct.create Init.sizeof in
    let (_: Cstruct.t) = Init.marshal Init.default buf in
    error_of_failure
      (fun () ->
         Lwt_cstruct.(complete (write fd) buf)
         >>= fun () ->
         Lwt_cstruct.(complete (read fd) buf)
         >>= fun () ->
         let open Infix in
         Lwt.return (Init.unmarshal buf)
         >>= fun (init, _) ->
         Lwt.return (`Ok { fd })
      )

  let simple_bool_command cmd t =
    error_of_failure
      (fun () ->
         let buf = Cstruct.create Command.sizeof in
         let (_: Cstruct.t) = Command.marshal cmd buf in
         Lwt_cstruct.(complete (write t.fd) buf)
         >>= fun () ->
         let result = Cstruct.create 1 in
         Lwt_cstruct.(complete (read t.fd) result)
         >>= fun () ->
         match Cstruct.get_uint8 result 0 with
         | 0 -> Lwt.return (`Ok ())
         | n -> Lwt.return (`Error (`Msg (Printf.sprintf "Command failed with code %d" n)))
      )

  let install_symlinks = simple_bool_command Command.Install_symlinks

end

let t =
  let open Lwt in
  let rec loop () =
    let s = Lwt_unix.socket Lwt_unix.PF_UNIX Lwt_unix.SOCK_STREAM 0 in
    Lwt.catch
      (fun () ->
        Lwt_unix.connect s (Lwt_unix.ADDR_UNIX "/var/tmp/com.docker.vmnetd.socket")
        >>= fun () ->
        Lwt.return s
      ) (fun e ->
        Lwt_unix.close s
        >>= fun () ->
        Lwt_unix.sleep 0.1
        >>= fun () ->
        loop ()
      ) in
  loop ()
  >>= fun s ->
  Vmnet_client.of_fd s
  >>= function
  | `Error (`Msg m) -> failwith m
  | `Ok t ->
    begin Vmnet_client.install_symlinks t
    >>= function
    | `Error (`Msg m) -> failwith m
    | `Ok () -> Lwt_unix.close s
    end

let _ = Lwt_main.run t
