(*
 * Copyright (c) 2016 Docker Inc. All rights reserved.
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 *)

module Log = Log.Log

type packet = {
  typ : char;
  ptr : unit Ctypes.ptr;
  size : int;
  now : int64;
}

type channel = {
  id : int;
  sink : (packet option -> unit) option;
  read : int -> Unsigned.uint8 Ctypes.CArray.t Lwt.t;
  write : Unsigned.uint8 Ctypes.ptr -> int -> int Lwt.t;
  mount : Mount.t;
}

type session = {
  k : int;
  name : string;
}

type t = {
  mutable session : session;
  mutable channels : channel option array;
}

let session_name = Printf.sprintf "osxfs.%04d.trace"

let new_tracer () =
  let k = 0 in
  { session = { k; name = session_name k}; channels = [||]; }

let get_trace_dir_name tracer =
  let log_dir = Osx_reporter.get_trace_dir () in
  let k = tracer.session.k in
  let name = session_name k in
  Filename.concat log_dir name, name

let rec create_session tracer =
  let trace_dir, name = get_trace_dir_name tracer in
  try
    Unix.mkdir trace_dir 0o700;
    tracer.session <- { tracer.session with name };
    trace_dir
  with
  | Unix.Unix_error (Unix.EEXIST, "mkdir", _) ->
    tracer.session <- { tracer.session with k = tracer.session.k + 1 };
    create_session tracer

let trace_head = Bytes.of_string "FUSESL\007\023" (* TODO: dynamic *)
let trace_head_len = Bytes.length trace_head

let le_bytes_of_int64 i =
  let byte_of_int64 k =
    char_of_int Int64.(to_int (logand (shift_right_logical i (8*k)) 0xff_L))
  in
  let buf = Bytes.create 8 in
  Bytes.set buf 0 (byte_of_int64 0);
  Bytes.set buf 1 (byte_of_int64 1);
  Bytes.set buf 2 (byte_of_int64 2);
  Bytes.set buf 3 (byte_of_int64 3);
  Bytes.set buf 4 (byte_of_int64 4);
  Bytes.set buf 5 (byte_of_int64 5);
  Bytes.set buf 6 (byte_of_int64 6);
  Bytes.set buf 7 (byte_of_int64 7);
  buf

let write_packet fd { typ; ptr; size; now } =
  let open Lwt.Infix in
  let fail = Lwt.fail (Failure "writing trace packet failed") in
  let lwt_fd = Lwt_unix.of_unix_file_descr fd in
  Lwt_unix.write lwt_fd (Bytes.make 1 typ) 0 1
  >>= fun typ_count ->
  if typ_count <> 1
  then (
    Log.err (fun f ->
      f "osxfs.trace: error writing packet type, wrote only %d bytes" typ_count
    );
    fail
  )
  else
    Lwt_unix.write lwt_fd (le_bytes_of_int64 now) 0 8
    >>= fun time_count ->
    if time_count <> 8
    then (
      Log.err (fun f ->
        f "Error writing timestamp, wrote only %d bytes" time_count
      );
      fail
    )
    else
      Unistd_unix_lwt.write fd ptr size
      >>= fun written ->
      if written <> size
      then (
        Log.err (fun f ->
          f "Error writing %d byte packet, wrote %d bytes" size written
        );
        fail
      )
      else Lwt.return_unit

let start_tracing tracer =
  let trace_dir = create_session tracer in
  Array.iteri (fun i -> function
    | None -> ()
    | Some channel ->
      let name = Mount.segment_name channel.mount in
      let filename = Filename.concat trace_dir ("trace"^name^".session") in
      let fd = Unix.(openfile filename [O_WRONLY; O_CREAT; O_EXCL] 0o600) in
      let head_count = Unix.write fd trace_head 0 trace_head_len in
      (if head_count <> trace_head_len
       then Log.err (fun f ->
         f "Error writing trace header, wrote only %d bytes" head_count
       ));
      let id = channel.id in
      let stream, push = Lwt_stream.create () in
      Lwt.async (fun () ->
        let open Lwt.Infix in
        let rec writer () =
          Lwt.catch (fun () ->
            Lwt_stream.get stream
            >>= function
            | None -> Lwt_unix.close (Lwt_unix.of_unix_file_descr fd)
            | Some packet ->
              write_packet fd packet
              >>= writer
          ) (fun exn ->
            Lwt_unix.close (Lwt_unix.of_unix_file_descr fd)
            >>= fun () ->
            Lwt.fail exn
          )
        in
        writer ()
      );
      match channel.sink with
      | Some sink ->
        sink None;
        tracer.channels.(i) <- Some { channel with sink = Some push }
      | None ->
        let current_socket = Fuse_lwt.get_socket id in
        Fuse_lwt.(set_socket id ~read:channel.read ~write:channel.write ());
        let read = Fuse_lwt.read_socket current_socket in
        let write = Fuse_lwt.write_socket current_socket in
        tracer.channels.(i) <- Some {
          channel with sink = Some push; read; write;
        }
  ) tracer.channels;
  trace_dir

let stop_tracing tracer =
  let trace_dir, name = get_trace_dir_name tracer in
  Array.iteri (fun i -> function
    | None -> ()
    | Some channel ->
      match channel.sink with
      | Some sink ->
        sink None;
        let id = channel.id in
        let current_socket = Fuse_lwt.get_socket id in
        Fuse_lwt.(set_socket id ~read:channel.read ~write:channel.write ());
        let read = Fuse_lwt.read_socket current_socket in
        let write = Fuse_lwt.write_socket current_socket in
        tracer.channels.(i) <- Some { channel with sink = None; read; write }
      | None -> ()
  ) tracer.channels;
  trace_dir

let push_packet typ tracer id ptr size now =
  match tracer.channels.(id) with
  | None -> Log.err (fun f -> f "osxfs.trace: trace channel %d not found" id)
  | Some { sink = None } ->
    Log.err (fun f -> f "osxfs.trace: sink for channel %d not found" id)
  | Some { sink = Some sink } ->
    sink (Some { typ; ptr; size; now })

let write_request = push_packet 'Q'
let write_reply = push_packet 'R'

let add_channel tracer ~id mount =
  let open Lwt.Infix in
  let base_socket = Fuse_lwt.get_socket id in
  let read size =
    Fuse_lwt.read_socket base_socket size
    >>= fun buffer ->
    let now = Mtime.(to_ns_uint64 (absolute ())) in
    let ptr = Ctypes.(to_voidp (CArray.start buffer)) in
    let size = Ctypes.CArray.length buffer in
    write_request tracer id ptr size now;
    Lwt.return buffer
  in
  let write ptr size =
    Fuse_lwt.write_socket base_socket ptr size
    >>= fun written ->
    let now = Mtime.(to_ns_uint64 (absolute ())) in
    write_reply tracer id (Ctypes.to_voidp ptr) written now;
    Lwt.return written
  in
  let table_size = Array.length tracer.channels in
  begin if id > table_size - 1 then
      let table = Array.make (id * 2 + 1) None in
      Array.blit tracer.channels 0 table 0 table_size;
      tracer.channels <- table
  end;
  let channel = { id; sink = None; read; write; mount; } in
  try
    tracer.channels.(id) <- Some channel
  with exn ->
    Log.err (fun f -> f "boom %s" (Printexc.to_string exn));
    raise exn
