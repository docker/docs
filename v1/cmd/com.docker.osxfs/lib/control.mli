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

type t

type log_level =
  | Error
  | Notice

module type HANDLER = sig
  val init : t -> unit Lwt.t
  val log : t -> log_level -> string -> unit Lwt.t
  val pong : t -> unit Lwt.t
  val set_notify_channel : Unix.file_descr -> string -> unit Lwt.t
end

val make : Unix.file_descr -> Lwt_unix.sockaddr -> t

val continue : t -> unit Lwt.t

val partition_suitable_mounts : t -> string list -> (string list * string list) Lwt.t

val partition_suitable_exports : t -> string list -> (string list * string list) Lwt.t

val service_connection : (module HANDLER) -> t -> unit -> unit Lwt.t
