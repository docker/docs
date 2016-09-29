open Lwt.Infix

module Log = Log.Log

type pending =
  | Suitable_mounts of string list * string list Lwt_mvar.t

type t = {
  client_fd : Unix.file_descr;
  client_addr : Lwt_unix.sockaddr;
  pending : (int, pending) Hashtbl.t;
  mutable max : int;
  mutable free : int list;
}

type log_level =
  | Error
  | Notice

let continue_message = 0
let mount_suitability_message = 1
let export_suitability_message = 2

module type HANDLER = sig
  val init : t -> unit Lwt.t
  val log : t -> log_level -> string -> unit Lwt.t
  val pong : t -> unit Lwt.t
  val set_notify_channel : Unix.file_descr -> string -> unit Lwt.t
end

let make client_fd client_addr = {
  client_fd;
  client_addr;
  pending = Hashtbl.create 8;
  max = 0;
  free = [];
}

let txn_id ctl = match ctl with
  | { free = []; max } -> let id = ctl.max in ctl.max <- id + 1; id
  | { free = id::free; } -> ctl.free <- free; id

let continue { client_fd } =
  Log.info (fun f -> f "sending continue to client");
  let len = 6 in
  let buf = Cstruct.create len in
  Cstruct.LE.set_uint32 buf 0 (Int32.of_int len);
  Cstruct.LE.set_uint16 buf 4 continue_message;
  Socket.write_exactly client_fd
    Ctypes.(bigarray_start array1 (Cstruct.to_bigarray buf)) len

let partition_suitable_paths msg_type ctl paths =
  let count = List.length paths in
  let hlen = 4 + 2 + 2 in
  let len =
    List.fold_left (fun acc s ->
      acc + String.length s
    ) (hlen + (count * 3)) paths
  in
  let buf = Cstruct.create len in
  Cstruct.LE.set_uint32 buf 0 (Int32.of_int len);
  Cstruct.LE.set_uint16 buf 4 msg_type;
  let id = txn_id ctl in
  let mvar = Lwt_mvar.create_empty () in
  Hashtbl.replace ctl.pending id (Suitable_mounts (paths, mvar));
  Cstruct.LE.set_uint16 buf 6 id;
  let _off = List.fold_left (fun off path ->
    let plen = String.length path in
    Cstruct.LE.set_uint16 buf off plen;
    Cstruct.blit_from_string path 0 buf (off + 2) plen;
    Cstruct.set_uint8 buf (off + 2 + plen) 0;
    off + plen + 3
  ) hlen paths in
  Socket.write_exactly ctl.client_fd
    Ctypes.(bigarray_start array1 (Cstruct.to_bigarray buf)) len
  >>= fun () ->
  Lwt_mvar.take mvar
  >|= fun suitable_mounts ->
  List.partition (fun path -> List.mem path suitable_mounts) paths

let partition_suitable_mounts =
  partition_suitable_paths mount_suitability_message

let partition_suitable_exports =
  partition_suitable_paths export_suitability_message

let respond_string_list ctl id strings =
  ctl.free <- id :: ctl.free;
  match Hashtbl.find ctl.pending id with
  | exception Not_found ->
    Lwt.fail (Failure "control transaction id not found")
  | Suitable_mounts (_, mvar) ->
    Lwt_mvar.put mvar strings
    >>= fun () ->
    Hashtbl.remove ctl.pending id;
    Lwt.return_unit

let rec service_connection handler ctl () =
  let module Handler = (val handler : HANDLER) in
  Socket.read_message ctl.client_fd
  >>= fun carray ->
  let sa = Ctypes.(coerce (ptr uint8_t) (ptr char) (CArray.start carray)) in
  let ca = Ctypes.CArray.(from_ptr sa (length carray)) in
  let ba = Ctypes.(bigarray_of_array array1 Bigarray.Char) ca in
  let cs = Cstruct.of_bigarray ba in
  let body = Cstruct.shift cs 4 in
  match Cstruct.LE.get_uint16 body 0 with
  | 0 ->
    Log.debug (fun f -> f "received init from client");
    Lwt.async (fun () -> Handler.init ctl);
    service_connection handler ctl ()
  | 1 ->
    let msg = Cstruct.(to_string (shift body 2)) in
    Lwt.async (fun () -> Handler.log ctl Error msg);
    service_connection handler ctl ()
  | 2 ->
    let msg = Cstruct.(to_string (shift body 2)) in
    Lwt.async (fun () -> Handler.log ctl Notice msg);
    service_connection handler ctl ()
  | 3 ->
    Lwt.async (fun () -> Handler.pong ctl);
    service_connection handler ctl ()
  | 4 ->
    let msg = Cstruct.shift body 2 in
    let id = Cstruct.LE.get_uint16 msg 0 in
    let msg = Cstruct.shift msg 2 in
    let rec get_paths paths cs =
      if Cstruct.len cs = 0
      then paths
      else
        let len = Cstruct.LE.get_uint16 cs 0 in
        let path = Cstruct.copy cs 2 len in
        get_paths (path::paths) (Cstruct.shift cs (len + 3))
    in
    let paths = get_paths [] msg in
    respond_string_list ctl id paths
    >>= fun () ->
    service_connection handler ctl ()
  | 5 ->
    let mount_point = Cstruct.(to_string (shift body 2)) in
    Handler.set_notify_channel ctl.client_fd mount_point
  | x -> Log.err (fun f -> f "Unknown transfuse request code %d" x); exit 1
