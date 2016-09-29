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

val inotify_channel : string -> (Lwt_unix.file_descr * Unsigned.UInt32.t) Lwt.t

module type CONF = sig
  val max_write : int
  val connect_path : string
  val address : string
end

module Mount(Conf : CONF)(FS:Fuse_lwt.FS_LWT)(IO: Fuse_lwt.IO_LWT)
  : Fuse.MOUNT_IO
    with module IO = IO
     and type t = FS.t

val listen_forever : string -> (module Control.HANDLER) -> unit Lwt.t
