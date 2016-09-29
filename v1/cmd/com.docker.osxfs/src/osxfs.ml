(*
 * Copyright (c) 2016 Docker Inc. All rights reserved.
 * Copyright (c) 2014-2015 David Sheets <sheets@alum.mit.edu>
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

module Profuse = Profuse_7_23

module LogInfo = Log.LogInfo
module Log = Log.Log

open Unsigned

module Stat = Sys_stat

module OsxfsData = struct
  type t = {
    path : string list;
  }

  let to_path { path } = path
  let create_child parent name = { path = name::parent.path }
  let rename _data path = { path }
end

module Node = Nodes.Path(OsxfsData)

module N = Nodes.Make(Node)

type state = {
  mount          : Mount.t;
  nodes          : N.t;
  max_read       : int;
}
type t = state

let max_write = 1 lsl 20
let max_read = 1 lsl 20
let readdir_block_count = 16

(* See branch osxfs-dcache-experiment for tools to help decide how to
   set this value. *)
let entry_valid_nsec = 10_000_000_l (* 10ms = 1 Linux jiffy *)
let entry_valid_nsec = UInt32.of_int32 entry_valid_nsec

let ping, ping_push = Lwt_stream.create ()
let ping_reply, ping_reply_push = Lwt_stream.create ()
let event_sock = ref Lwt_unix.stdout

(* TODO: check N.get raising Not_found *)
(* TODO: check uid/gid rights *)
(* TODO: check fuse_reply_{entry,create} invocations increment lookup *)

(* TODO: this lock table should live in the node table's parameter *)
type fh_state =
  | Closed
  | Open

type fh_resource = {
  fh : int64;
  mutex : Lwt_mutex.t;
  mutable state : fh_state;
}

let fh_locks = Hashtbl.create 256

let get_fh_resource fh =
  try
    Hashtbl.find fh_locks fh
  with Not_found ->
    let mutex = Lwt_mutex.create () in
    let resource = { fh; mutex; state = Open } in
    Hashtbl.replace fh_locks fh resource;
    resource

let unlock_fh_resource rsrc state =
  rsrc.state <- state;
  if state = Closed then Hashtbl.remove fh_locks rsrc.fh;
  Lwt_mutex.unlock rsrc.mutex

let is_uncontended { mutex } = Lwt_mutex.is_empty mutex

let with_fh_resource fh f =
  let open Lwt.Infix in
  let { mutex } as resource = get_fh_resource fh in
  Lwt_mutex.lock mutex
  >>= fun () ->
  let state = resource.state in
  Lwt.catch (fun () ->
    f state
    >>= fun (state, v) ->
    unlock_fh_resource resource state;
    Lwt.return v
  ) (fun exn ->
    let state = match exn with
      | Unix.Unix_error (Unix.EBADF, _, _) -> Closed
      | Errno.Error { Errno.errno } when List.mem Errno.EBADF errno -> Closed
      | _other -> state
    in
    unlock_fh_resource resource state;
    Lwt.fail exn
  )

let disable_dcache_table = Hashtbl.create 64

let disable_dcache path =
  (* We use 'add' in order to track the number of in-flight events. *)
  Hashtbl.add disable_dcache_table path ()

let enable_dcache parent name_opt =
  let ppath = Node.to_string parent.Nodes.data in
  let path = match name_opt with
    | None -> ppath
    | Some name -> Filename.concat ppath name
  in
  if Hashtbl.mem disable_dcache_table path
  then Hashtbl.remove disable_dcache_table path
  else Log.err (fun f ->
    f "INVARIANT VIOLATED: DCACHE ALREADY ENABLED FOR %s" path
  )

let is_path_dcached path =
  not (Hashtbl.mem disable_dcache_table path)

let is_node_dcached node = is_path_dcached (Node.to_string node.Nodes.data)

let string_of_nodeid nodeid st =
  N.string_of_id st.nodes (UInt64.to_int64 nodeid)

let uint64_of_int64 = UInt64.of_int64
let uint32_of_uint64 x = UInt32.of_int (UInt64.to_int x)

let make mount root =
  let path = match root with
    | "" -> [""]
    | root -> Stringext.split root ~on:'/'
  in
  let label = "nodes_"^root in
  let root = Node.of_value { OsxfsData.path = List.rev path; } in
  {
    mount;
    nodes    = N.create ~label root;
    max_read;
  }

let log_error message = Log.err (fun f -> f "%s" message)

module FS : Fuse_lwt.FS_LWT with type t = state = struct
  type t = state

  let string_of_nodeid = string_of_nodeid

  let string_of_state req st =
    Printf.sprintf "Nodes: %s" (N.to_string st.nodes)

  let log_error = log_error

  module Calls(IO : Fuse_lwt.IO_LWT)
    : Fuse_lwt.FS_IO_LWT with type 'a IO.t = 'a IO.t and type t = t = struct
    module IO = IO
    type 'a req_handler = Fuse.request -> 'a -> 'a IO.t

    type t = state

    open Lwt

    module Support = Fuse.Support(IO)

    (* will be overridden *)
    let dispatch _req st = return st

    let negotiate_mount pkt req st =
      ignore (Unix.umask 0o000);
      let open Unsigned in
      let open Profuse in
      let major = UInt32.to_int (Ctypes.getf pkt In.Init.T.major) in
      if major <> 7
      then IO.fail (
        ProtocolError
          (req.chan, Printf.sprintf
             "Incompatible FUSE protocol major version %d <> 7" major))
      else
        let minor = UInt32.to_int (Ctypes.getf pkt In.Init.T.minor) in
        if minor < 8
        then IO.fail (
          ProtocolError
            (req.chan, Printf.sprintf
               "Incompatible FUSE protocol minor version %d < 8" minor))
        else
          let minor = min minor 23 in (* TODO: track kernel minor *)
          let max_readahead = Ctypes.getf pkt In.Init.T.max_readahead in
          let max_readahead = min (UInt32.to_int max_readahead) st.max_read in
          (* TODO: ? *)
          let _flags = Ctypes.getf pkt In.Init.T.flags in (* TODO: ? *)
          let chan = {
            req.chan with version = (major, minor); max_readahead; flags = 0l;
          } in
          let major = UInt32.of_int major in
          let minor = UInt32.of_int minor in
          Log.debug (fun f -> f "max_readahead = %d max_write = %d"
            max_readahead req.chan.max_write);
          let max_readahead = UInt32.of_int max_readahead in
          let max_write = UInt32.of_int req.chan.max_write in
          let flags = UInt32.Infix.(Types.In.Init.Flags.fuse_big_writes
                                lor Types.In.Init.Flags.async_read) in
          let pkt =
            Out.Init.create ~major ~minor ~max_readahead ~flags ~max_write
          in
          IO.(Out.write_reply req pkt
              >>= fun () ->
              return ({req with chan}, st)
             )

    let enosys = Support.enosys
    let nodeid req = UInt64.to_int64 (Support.nodeid req)

    let flags_of_code req code =
      let phost = Profuse.(req.chan.host.Host.fcntl.Fcntl.Host.oflags) in
      List.fold_left Fcntl.Oflags.(fun list -> function
        | O_NOATIME (* TODO: emulate? *)
        | O_RSYNC ->
          (* Linux aliases O_SYNC and O_RSYNC but OS X only has O_SYNC *)
          list
        | flag ->
          let host = Fcntl_unix.Oflags.host in
          match Fcntl.Oflags.to_code ~host [flag] with
          | exception Not_found -> failwith "NO_OFLAG"
          | _code -> flag :: list
      ) [] (Fcntl.Oflags.of_code ~host:phost (UInt32.to_int code))

    (* TODO: caching + events could potentially improve this *)
    let get_ownership f path =
      Lwt.catch (fun () -> f path)
        (function
          | Errno.Error { Errno.errno } when List.mem Errno.EACCES errno ->
            Lwt.return Xowner.empty
          | exn -> Lwt.fail exn
        )

    let get_fd_fh_of_node node =
      let rec next p = Nodes.UnixHandle.(function
        | [] -> None
        | (_, Dir (_,_)) :: rest -> next p rest
        | (fh, File (fd,_)) :: rest ->
          let resource = get_fh_resource fh in
          if p resource
          then Some (fd, fh)
          else next p rest
      ) in
      let handles = N.handles node in
      match next is_uncontended handles with
      | None -> next (fun _ -> true) handles
      | some_fd_resource -> some_fd_resource

    let store_attr_of_stat_owner req (stat_, owner) = Lwt_unix.LargeFile.(
      let req_uid = Ctypes.getf req.Profuse.hdr Profuse.In.Hdr.T.uid in
      let req_gid = Ctypes.getf req.Profuse.hdr Profuse.In.Hdr.T.gid in
      let uid = match owner.Xowner.uid with
        | None -> req_uid
        | Some id -> UInt32.of_int id
      in
      let gid = match owner.Xowner.gid with
        | None -> req_gid
        | Some id -> UInt32.of_int id
      in
      let open Sys_stat_unix in
      (* TODO: why does ino have to go to int first? *)
      let ino = PosixTypes.Ino.to_int (Stat.ino stat_) in
      let size = PosixTypes.Off.to_int64 (Stat.size stat_) in
      let blocks = Posix_types.Blkcnt.to_int64 (Stat.blocks stat_) in
      let blksize = Posix_types.Blksize.to_int64 (Stat.blksize stat_) in
      let atime = Posix_types.Time.to_int64 (Stat.atime stat_) in
      let mtime = Posix_types.Time.to_int64 (Stat.mtime stat_) in
      let ctime = Posix_types.Time.to_int64 (Stat.ctime stat_) in
      (* TODO: why does mode have to go to int first? *)
      let mode = Posix_types.Mode.to_int (Stat.mode stat_) in
      (* TODO: why does nlink have to go to int first? *)
      let nlink = Posix_types.Nlink.to_int (Stat.nlink stat_) in
      (* TODO: why does rdev have to go to int first? *)
      let rdev = Posix_types.Dev.to_int (Stat.rdev stat_) in
      return (Profuse.Struct.Attr.store
                ~ino:(UInt64.of_int ino)
                ~size:(uint64_of_int64 size)
                ~blocks:(uint64_of_int64 blocks)
                ~blksize:(UInt32.of_int64 blksize)
                ~atime:(uint64_of_int64 atime)
                ~atimensec:UInt32.zero
                ~mtime:(uint64_of_int64 mtime)
                ~mtimensec:UInt32.zero
                ~ctime:(uint64_of_int64 ctime)
                ~ctimensec:UInt32.zero
                ~mode:(UInt32.of_int mode)
                ~nlink:(UInt32.of_int nlink)
                ~uid
                ~gid
                ~rdev:(UInt32.of_int rdev)
             )
    )

    let rec store_attr_of_node req node =
      (* TODO: fix race access between lstat and Xowner.get *)
      match get_fd_fh_of_node node with
      | Some (fd, fh) -> with_fh_resource fh (function
        | Closed -> Lwt.return (Closed, None)
        | Open ->
          Sys_stat_unix_lwt.fstat fd
          >>= fun stat ->
          get_ownership Xowner.fget fd
          >>= fun owner ->
          Lwt.return (Open, Some (stat, owner))
      ) >>= (function
        | None -> store_attr_of_node req node
        | Some stat_owner -> store_attr_of_stat_owner req stat_owner
      )
      | None -> match node.Nodes.parent with
        | None -> Lwt.fail Unix.(Unix_error (ENOENT, "", ""))
        | Some _ ->
          let path_string = Node.to_string node.Nodes.data in
          Sys_stat_unix_lwt.lstat path_string
          >>= fun stat ->
          get_ownership Xowner.get path_string
          >>= fun owner ->
          store_attr_of_stat_owner req (stat, owner)

    let getattr req st =
      catch (fun () ->
        let node = N.get st.nodes (nodeid req) in
        let attr_valid = UInt64.zero in
        let attr_valid_nsec = UInt32.zero in
        store_attr_of_node req node
        >>= fun store_attr ->
        IO.Out.write_reply req
          (Profuse.Out.Attr.create ~attr_valid ~attr_valid_nsec ~store_attr)
        >>= fun () ->
        return st
      ) (function
        | Not_found ->
          Log.warn (fun f -> f "getattr not found");
          IO.Out.write_error log_error req Errno.ENOENT
          >>= fun () ->
          return st
        | e -> fail e
      )

    let opendir op req st =
      catch (fun () ->
        let { nodes } = st in
        let id = nodeid req in
        let node = N.get nodes id in
        let path = Node.to_string node.Nodes.data in
        Lwt_unix.opendir path
        >>= fun dir ->
        let h = N.Handle.open_ nodes id (Nodes.UnixHandle.Dir (dir, 0)) in
        let fh = UInt64.of_int64 h in
        let open_flags = Profuse.Out.Open.Flags.zero in
        IO.Out.write_reply req (Profuse.Out.Open.create ~fh ~open_flags)
        >>= fun () ->
        return st
      ) (function
        | Not_found ->
          Log.warn (fun f -> f "opendir not found");
          IO.Out.write_error log_error req Errno.ENOENT
          >>= fun () ->
          return st
        | e -> fail e
      )

    let forget_node id n st =
      try
        N.forget st.nodes id (UInt64.to_int n);
        st
      with Not_found ->
        Log.warn (fun f -> f "forget not found for node %Ld" id);
        st

    let forget n req st =
      (* FORGET is a non-returning command. *)
      return (forget_node (nodeid req) n st)

    let batch_forget forgets req st =
      let st = List.fold_left (fun st forget ->
        let id = Ctypes.getf forget Profuse.Struct.Forget_one.T.nodeid in
        let n = Ctypes.getf forget Profuse.Struct.Forget_one.T.nlookup in
        forget_node (UInt64.to_int64 id) n st
      ) st forgets in
      (* BATCH_FORGET is a non-returning command. *)
      return st

    let store_entry req node =
      let entry_valid_nsec =
        if is_node_dcached node
        then entry_valid_nsec
        else UInt32.zero
      in
      Support.store_entry ~entry_valid_nsec (store_attr_of_node req) node
    let respond_with_entry req node =
      let store_attr = store_attr_of_node req in
      let entry_valid_nsec =
        if is_node_dcached node
        then entry_valid_nsec
        else UInt32.zero
      in
      Support.respond_with_entry ~entry_valid_nsec store_attr node

    let lookup name req st =
      let parent = N.get st.nodes (nodeid req) in
      respond_with_entry req (N.lookup parent name) req
      >>= fun () ->
      return st

    let readdir r req st =
      let req_off = UInt64.to_int (Ctypes.getf r Profuse.In.Read.T.offset) in
      let rec seek dir off =
        if off < req_off
        then
          Dirent_unix_lwt.readdir dir
          >>= fun _ ->
          seek dir (off + 1)
        else if off > req_off
        then
          Lwt_unix.rewinddir dir
          >>= fun () ->
          seek dir 0
        else Lwt.return off
      in
      let host = Profuse.(req.chan.host) in
      let sz = UInt32.to_int (Ctypes.getf r Profuse.In.Read.T.size) in
      let fh = UInt64.to_int64 (Ctypes.getf r Profuse.In.Read.T.fh) in
      (match N.Handle.get st.nodes fh with
       | Nodes.UnixHandle.Dir (dir, off) ->
         with_fh_resource fh (function
           | Closed -> Lwt.fail Unix.(Unix_error (EBADF, "", ""))
           | Open ->
             seek dir off
             >>= fun off ->
             assert (off = req_off);
             let block_count = readdir_block_count in
             let rec readdir_k listing acc = function
               | 0 -> Lwt.return listing
               | k ->
                 Lwt.catch
                   (fun () ->
                      Dirent_unix_lwt.readdir dir
                      >|= fun dirent -> Some dirent
                   ) (function
                     | End_of_file -> Lwt.return_none
                     | exn -> Lwt.fail exn
                   )
                 >>= function
                 | Some { Dirent.Dirent.name; kind; ino } ->
                   let off = off + block_count + 1 - k in
                   N.Handle.set st.nodes fh (Nodes.UnixHandle.Dir (dir, off));
                   let acc = acc + Profuse.Out.Dirent.size name in
                   if acc > sz
                 then Lwt.return listing
                 else readdir_k ((off, ino, name, kind)::listing) acc (k - 1)
                 | None -> Lwt.return listing
             in
             readdir_k [] 0 block_count >|= fun listing -> Open, listing
         ) >>= fun dirents ->
         let listing = List.rev dirents in
         IO.Out.write_reply req
           (Profuse.Out.Dirent.of_list ~host listing 0 sz)
       | Nodes.UnixHandle.File (_,_) ->
         Lwt.fail Unix.(Unix_error (ENOTDIR, "", ""))
      ) >>= fun () ->
      return st

    (* Can raise Unix.Unix_error *)
    let readlink req st =
      let node = N.get st.nodes (nodeid req) in
      (* errors caught by our caller *)
      let path = Node.to_string node.Nodes.data in
      Lwt_unix.readlink path
      >>= fun target ->
      IO.Out.write_reply req (Profuse.Out.Readlink.create ~target)
      >>= fun () ->
      return st

    let open_ op req st =
      let open Profuse in
      catch (fun () ->
        let { nodes } = st in
        let id = nodeid req in
        let node = N.get nodes id in
        let path = Node.to_string node.Nodes.data in
        let flags = Ctypes.getf op In.Open.T.flags in
        let flags = flags_of_code req flags in
        Fcntl_unix_lwt.open_ path flags (* TODO: Is no perm OK? *)
        >>= fun fd ->
        Lwt_unix.(fstat (of_unix_file_descr fd))
        >>= fun { Unix.st_kind } ->
        let kind = Sys_stat_unix.File_kind.of_unix st_kind in
        let h = N.Handle.open_ nodes id (Nodes.UnixHandle.File (fd, kind)) in
        let fh = UInt64.of_int64 h in
        let open_flags = Profuse.Out.Open.Flags.zero in
        IO.Out.write_reply req (Out.Open.create ~fh ~open_flags)
        (* TODO: flags *)
        >>= fun () ->
        return st
      ) (function
        | Failure "NO_OFLAG" ->
          IO.Out.write_error log_error req Errno.EINVAL
          >>= fun () ->
          return st
        | Not_found ->
          (* TODO: log? *)
          IO.Out.write_error log_error req Errno.ENOENT
          >>= fun () ->
          return st
        | e -> fail e
      )

    let read r req st =
      let fh = UInt64.to_int64 (Ctypes.getf r Profuse.In.Read.T.fh) in
      let offset = Ctypes.getf r Profuse.In.Read.T.offset in
      let size = UInt32.to_int (Ctypes.getf r Profuse.In.Read.T.size) in
      let size = min size st.max_read in
      (match N.Handle.get st.nodes fh with
       | Nodes.UnixHandle.File (fd, _kind) ->
         let offset = UInt64.to_int64 offset in
         let pkt = Profuse.Out.Read.allocate ~size req in
         let ptr = Ctypes.to_voidp (Ctypes.CArray.start pkt) in
         with_fh_resource fh (function
           | Closed -> Lwt.fail Unix.(Unix_error (EBADF, "", ""))
           | Open ->
             Unistd_unix_lwt.pread fd ptr size offset
             >|= fun size -> Open, size
         ) >>= fun size ->
         IO.Out.write_reply req (Profuse.Out.Read.finalize ~size pkt)
       | Nodes.UnixHandle.Dir (_, _) ->
         Lwt.fail Unix.(Unix_error (EISDIR, "", ""))
      ) >>= fun () ->
      return st

    (* TODO: anything? *)
    let flush f req st = IO.Out.write_ack req >>= fun () -> return st

    (* TODO: flags? *)
    let release r req st =
      catch (fun () ->
        let fh = UInt64.to_int64 (Ctypes.getf r Profuse.In.Release.T.fh) in
        with_fh_resource fh (function
          | Closed -> Lwt.fail Unix.(Unix_error (EBADF, "", ""))
          | Open ->
            N.Handle.free st.nodes fh;
            Lwt.return (Closed, ())
        ) >>= fun () ->
        IO.Out.write_ack req
        >>= fun () ->
        return st
      ) (function
        | Not_found ->
          IO.Out.write_error log_error req Errno.EBADF
          >>= fun () ->
          return st
        | e -> fail e
      )

    (* TODO: distinguish? *)
    let releasedir = release

    (* Can raise Unix.Unix_error *)
    let symlink name target req st =
      let ({ Nodes.data } as pnode) = N.get st.nodes (nodeid req) in
      let path = Filename.concat (Node.to_string data) name in
      (* errors caught by our caller *)
      Lwt_unix.symlink target path
      >>= fun () ->
      lookup name req st (* TODO: still increment lookups? *)

    (* Can raise Unix.Unix_error *)
    let rename r src dest req st =
      let parent = nodeid req in
      let olddir = N.get st.nodes parent in
      let path = Node.to_string olddir.Nodes.data in
      let newdir = UInt64.to_int64 (Ctypes.getf r Profuse.In.Rename.T.newdir) in
      let newdir = N.get st.nodes newdir in
      let newpath = Node.to_string newdir.Nodes.data in
      (* errors caught by our caller *)
      Lwt_unix.rename (Filename.concat path src) (Filename.concat newpath dest)
      >>= fun () ->
      N.rename olddir src newdir dest;
      IO.Out.write_ack req
      >>= fun () ->
      return st

    (* Can raise Unix.Unix_error *)
    let unlink name req st =
      let node = N.get st.nodes (nodeid req) in
      let path = Node.to_string node.Nodes.data in
      let path = Filename.concat path name in
      (* errors caught by our caller *)
      Lwt_unix.unlink path
      >>= fun () ->
      N.unlink node name;
      IO.Out.write_ack req
      >>= fun () ->
      return st

    (* Can raise Unix.Unix_error *)
    let rmdir name req st =
      let node = N.get st.nodes (nodeid req) in
      let path = Node.to_string node.Nodes.data in
      let path = Filename.concat path name in
      (* errors caught by our caller *)
      Lwt_unix.rmdir path
      >>= fun () ->
      N.unlink node name;
      IO.Out.write_ack req
      >>= fun () ->
      return st

    (* Can raise Errno.Error *)
    let statfs req st =
      let node = N.get st.nodes (nodeid req) in
      let path = Node.to_string node.Nodes.data in
      (* TODO: lwtize *)
      let statfs = Osx_mount.statfs path in
      (* Note: optimal transfer block size is called iosize on OS X
         and bsize on Linux (OS X has an *actual* block size called
         `bsize`). *)
      (* TODO: experiment with adjusting for vsock channel optimization *)
      let bsize = UInt32.of_int statfs.Osx_mount.Statfs.iosize in
      let frsize = UInt32.of_int statfs.Osx_mount.Statfs.bsize in
      IO.Out.write_reply req Osx_mount.Statfs.(
        Profuse.Out.Statfs.create
          ~blocks:(UInt64.of_int64 statfs.blocks)
          ~bfree:(UInt64.of_int64 statfs.bfree)
          ~bavail:(UInt64.of_int64 statfs.bavail)
          ~files:(UInt64.of_int64 statfs.files)
          ~ffree:(UInt64.of_int64 statfs.ffree)
          ~bsize
          ~namelen:(UInt32.of_int maxpathlen)
          ~frsize
      )
      >>= fun () ->
      return st

    (* TODO: do *)
    let fsync _f = enosys

    (* Can raise Unix.Unix_error *)
    (* TODO: write flags? *)
    let write w ptr req st =
      let fh = UInt64.to_int64 (Ctypes.getf w Profuse.In.Write.T.fh) in
      let offset = UInt64.to_int64 (Ctypes.getf w Profuse.In.Write.T.offset) in
      let size = UInt32.to_int (Ctypes.getf w Profuse.In.Write.T.size) in
      (match N.Handle.get st.nodes fh with
       | Nodes.UnixHandle.File (fd, _kind) ->
         let ptr = Ctypes.(coerce (ptr char) (ptr void)) ptr in
         with_fh_resource fh (function
           | Closed -> Lwt.fail Unix.(Unix_error (EBADF, "", ""))
           | Open ->
             (* errors caught by our caller *)
             Unistd_unix_lwt.pwrite fd ptr size offset
             >|= fun written -> Open, written
         ) >>= fun written ->
         let size = UInt32.of_int written in
         IO.Out.write_reply req (Profuse.Out.Write.create ~size)
       | Nodes.UnixHandle.Dir (_, _) ->
         Lwt.fail Unix.(Unix_error (EISDIR, "", ""))
      )
      >>= fun () ->
      return st

    (* Can raise Unix.Unix_error *)
    let link l name req st =
      let { Nodes.data } = N.get st.nodes (nodeid req) in
      let path = Filename.concat (Node.to_string data) name in
      let oldnodeid =
        UInt64.to_int64 (Ctypes.getf l Profuse.In.Link.T.oldnodeid)
      in
      let oldnode = N.get st.nodes oldnodeid in
      let oldpath = Node.to_string oldnode.Nodes.data in
      (* errors caught by our caller *)
      Lwt_unix.link oldpath path
      >>= fun () ->
      lookup name req st (* TODO: still increment lookups? *)

    (* TODO: do *)
    let getxattr _g = enosys

    (* TODO: do *)
    let setxattr _s = enosys

    (* TODO: do *)
    let listxattr _g = enosys

    (* TODO: do *)
    let removexattr _name = enosys

    let access a req st =
      let { Nodes.data } = N.get st.nodes (nodeid req) in
      let path = Node.to_string data in
      let code = UInt32.to_int (Ctypes.getf a Profuse.In.Access.T.mask) in
      let phost = Profuse.(req.chan.host.Host.unistd.Unistd.access) in
      let perms = Unistd.Access.(of_code ~host:phost code) in
      Lwt_unix.access path perms
      >>= fun () ->
      IO.Out.write_ack req
      >>= fun () ->
      return st

    let create c name req st =
      catch (fun () ->
        let { nodes } = st in
        let parent_id = nodeid req in
        let pnode = N.get nodes parent_id in
        let path = Node.to_string pnode.Nodes.data in
        let mode = UInt32.to_int (Ctypes.getf c Profuse.In.Create.T.mode) in
        (* TODO: is only file_perm? *)
        let flags = Ctypes.getf c Profuse.In.Create.T.flags in
        let flags = flags_of_code req flags in
        let path = Filename.concat path name in
        Fcntl_unix_lwt.open_ path ~perms:mode (Fcntl.Oflags.O_EXCL::flags)
        >>= fun fd ->
        Lwt_unix.(fstat (of_unix_file_descr fd))
        >>= fun stat ->
        let kind = Sys_stat_unix.File_kind.of_unix stat.Unix.st_kind in
        let node = N.lookup pnode name in
        let h = Nodes.UnixHandle.File (fd, kind) in
        let fh = N.Handle.open_ nodes node.Nodes.id h in
        store_entry req node
        >>= fun store_entry ->
        IO.Out.write_reply req
          (Profuse.Out.Create.create
             ~store_entry
             ~store_open:(
               let fh = UInt64.of_int64 fh in
               let open_flags = Profuse.Out.Open.Flags.zero in
               Profuse.Out.Open.store ~fh ~open_flags
             ))
        >>= fun () ->
        (* TODO: flags *)
        return st
      ) (function
        | Failure "NO_OFLAG" ->
          IO.Out.write_error log_error req Errno.EINVAL
          >>= fun () ->
          return st
        | Not_found ->
          (* TODO: log? *)
          IO.Out.write_error log_error req Errno.ENOENT
          >>= fun () ->
          return st
        | e -> fail e
      )

    let mknod m name req st =
      let ({ Nodes.data } as pnode) = N.get st.nodes (nodeid req) in
      let path = Node.to_string data in
      let path = Filename.concat path name in
      let phost = Profuse.(req.chan.host.Host.sys_stat.Sys_stat.Host.mode) in
      let mode_code =
        Unsigned.UInt32.to_int (Ctypes.getf m Profuse.In.Mknod.T.mode)
      in
      Stat.File_kind.(match Stat.Mode.of_code ~host:phost mode_code with
        | Some (FIFO, perm) ->
          Lwt_unix.mkfifo path perm
          >>= fun () ->
          respond_with_entry req (N.lookup pnode name) req
          >>= fun () ->
          return st
        | Some (SOCK, perm) ->
          (* We just want to create the node, we don't care about the
             transport. TODO: care about the transport *)
          let sock_fd = Lwt_unix.(socket PF_UNIX SOCK_STREAM 0) in
          let () = Lwt_unix.(bind sock_fd (ADDR_UNIX path)) in
          Lwt_unix.close sock_fd
          >>= fun () ->
          respond_with_entry req (N.lookup pnode name) req
          >>= fun () ->
          return st
        | Some ((DIR | CHR | BLK | REG | LNK), perm) ->
          enosys req st (* TODO: implement? different error? *)
        | None ->
          IO.Out.write_error log_error req Errno.ENOTSUP
          >>= fun () ->
          return st
      )
    (*
    let ({ Nodes.data } as pnode) = N.get st.nodes (nodeid req) in
    let path = Nodes.Path.to_string data in
    let path = Filename.concat path name in
    let mode = Ctypes.getf m Profuse.In.Mknod.T.mode in
    let rdev = Ctypes.getf m Profuse.In.Mknod.T.rdev in (* TODO: use this? *)
    (* TODO: translate mode and dev from client host rep to local host rep *)
    (* TODO: dev_t is usually 64-bit but rdev is 32-bit. translate how? *)
    (* TODO: regular -> open with O_CREAT | O_EXCL | O_WRONLY for compat? *)
    (* TODO: fifo -> mkfifo for compat? *)
    agents.Agent_handler.mknod ~uid ~gid path mode
      (Unsigned.UInt32.to_int32 rdev);
    respond_with_entry (N.lookup pnode name) req;
    st
*)

    let mkdir m name req st =
      let ({ Nodes.data } as pnode) = N.get st.nodes (nodeid req) in
      let path = Node.to_string data in
      let path = Filename.concat path name in
      let mode = UInt32.to_int (Ctypes.getf m Profuse.In.Mkdir.T.mode) in
      Lwt_unix.mkdir path mode
      >>= fun () ->
      respond_with_entry req (N.lookup pnode name) req
      >>= fun () ->
      return st

    (* TODO: do *)
    let fsyncdir _f = enosys

    (* TODO: do *)
    let getlk _lk = enosys

    (* TODO: do *)
    let setlk _lk = enosys

    (* TODO: do *)
    let setlkw _lk = enosys

    (* TODO: do *)
    (* The correct response for unimplemented interrupt support is to ignore *)
    let interrupt _i _req st = return st

    (* TODO: do *)
    let bmap _b = enosys

    let destroy _req st = return st

    let set_ownership ?uid ?gid fd =
      Xowner.fset ?uid ?gid fd

    (* TODO: mode could be considered when deciding how to order
       ownership and permissions changes *)
    let setattr_fd valid s req fd k = Profuse.In.Setattr.(
      let set_uid = valid.Valid.uid in
      let set_gid = valid.Valid.gid in
      let uid =
        if set_uid then Some (UInt32.to_int (Ctypes.getf s T.uid)) else None
      in
      let gid =
        if set_gid then Some (UInt32.to_int (Ctypes.getf s T.gid)) else None
      in
      if set_uid || set_gid
      then
        let fd = Lwt_unix.unix_file_descr fd in
        set_ownership ?uid ?gid fd
      else return_unit
      >>= fun () ->

      (if valid.Valid.mode
       then
         let mode = UInt32.to_int (Ctypes.getf s T.mode) in
         let phost = Profuse.(
           req.chan.host.Host.sys_stat.Sys_stat.Host.mode
         ) in
         let (kind,perm) = Stat.Mode.(
           of_code_exn ~host:phost mode
         ) in
         (if kind <> k
          then fail Errno.(Error {
            errno = [ EOPNOTSUPP ];
            call  = "setattr";
            label = "";
          })
          else return_unit)
         >>= fun () ->
         Lwt_unix.fchmod fd perm
       else return_unit)
      >>= fun () ->

      (if valid.Valid.size
       then
         let size = UInt64.to_int64 (Ctypes.getf s T.size) in
         Lwt_unix.LargeFile.ftruncate fd size
       else return_unit)
      >>= fun () ->

      (* TODO: support ns for host file systems that can do that (not HFS+) *)
      (if (valid.Valid.atime ||
           valid.Valid.atime_now ||
           valid.Valid.mtime ||
           valid.Valid.mtime_now)
       then
         let open Osx_attr in
         let now = Int64.of_float (Unix.gettimeofday ()) in
         let atime =
           if valid.Valid.atime
           then let atime = Ctypes.getf s T.atime in [
             Value.Common (Common.ACCTIME, {
               Time.Timespec.sec = Unsigned.UInt64.to_int64 atime;
               nsec = 0;
             }) ]
           else if valid.Valid.atime_now
           then [
             Value.Common (Common.ACCTIME, {
               Time.Timespec.sec = now;
               nsec = 0;
             }) ]
           else []
         in
         let mtime =
           if valid.Valid.mtime
           then let mtime = Ctypes.getf s T.mtime in [
             Value.Common (Common.MODTIME, {
               Time.Timespec.sec = Unsigned.UInt64.to_int64 mtime;
               nsec = 0;
             }) ]
           else if valid.Valid.mtime_now
           then [
             Value.Common (Common.MODTIME, {
               Time.Timespec.sec = now;
               nsec = 0;
             }) ]
           else []
         in
         let fd = Lwt_unix.unix_file_descr fd in
         Osx_attr_lwt.fsetlist ~no_follow:true (atime @ mtime) fd
       else Lwt.return_unit
      )
      (* We can't set ctime on OS X so we just ignore it. *)
    )

    let setattr_path valid s req path = Profuse.In.Setattr.(
      Lwt_unix.stat path
      >>= fun { Unix.st_kind } ->

      (* TODO: check other types *)
      (if valid.Valid.size && st_kind = Unix.S_SOCK
       then fail Unix.(Unix_error (EINVAL, "setattr", path))
       else (* we change the size later... *) return_unit)
      >>= fun () ->

      (if valid.Valid.mode
       then
         let mode = UInt32.to_int (Ctypes.getf s T.mode) in
         let phost = Profuse.(
           req.chan.host.Host.sys_stat.Sys_stat.Host.mode
         ) in
         let (kind,perm) = Stat.Mode.(
           of_code_exn ~host:phost mode
         ) in
         (if Sys_stat_unix.File_kind.to_unix kind <> st_kind
          then fail Errno.(Error {
            errno = [ EOPNOTSUPP ];
            call = "setattr";
            label = "";
          })
          else return_unit)
         >>= fun () ->
         Lwt_unix.chmod path perm
       else return_unit)

      >>= fun () ->

      let set_uid = valid.Valid.uid in
      let set_gid = valid.Valid.gid in
      (if set_uid
       then Xowner.set ~uid:(UInt32.to_int (Ctypes.getf s T.uid)) path
       else if set_gid
       then Xowner.set ~gid:(UInt32.to_int (Ctypes.getf s T.gid)) path
       else Lwt.return_unit
      ) >>= fun () ->

      (if valid.Valid.size
       then
         let size = UInt64.to_int64 (Ctypes.getf s T.size) in
         Lwt_unix.LargeFile.truncate path size
       else return_unit)
      >>= fun () ->

      (* TODO: support ns for host file systems that can do that (not HFS+) *)
      (if (valid.Valid.atime ||
           valid.Valid.atime_now ||
           valid.Valid.mtime ||
           valid.Valid.mtime_now)
       then
         let open Osx_attr in
         let now = Int64.of_float (Unix.gettimeofday ()) in
         let atime =
           if valid.Valid.atime
           then let atime = Ctypes.getf s T.atime in [
             Value.Common (Common.ACCTIME, {
               Time.Timespec.sec = Unsigned.UInt64.to_int64 atime;
               nsec = 0;
             }) ]
           else if valid.Valid.atime_now
           then [
             Value.Common (Common.ACCTIME, {
               Time.Timespec.sec = now;
               nsec = 0;
             }) ]
           else []
         in
         let mtime =
           if valid.Valid.mtime
           then let mtime = Ctypes.getf s T.mtime in [
             Value.Common (Common.MODTIME, {
               Time.Timespec.sec = Unsigned.UInt64.to_int64 mtime;
               nsec = 0;
             }) ]
           else if valid.Valid.mtime_now
           then [
             Value.Common (Common.MODTIME, {
               Time.Timespec.sec = now;
               nsec = 0;
             }) ]
           else []
         in
         Osx_attr_lwt.setlist ~no_follow:true (atime @ mtime) path
       else Lwt.return_unit
      )
      (* We can't set ctime on OS X so we just ignore it. *)
    )

    let setattr s req st = Profuse.In.Setattr.(
      let valid = Valid.of_uint32 (Ctypes.getf s T.valid) in
      if valid.Valid.fh
      then
        (* TODO: can fh be readonly and cause truncate problems?
           test with ftruncate
        *)
        let fh = UInt64.to_int64 (Ctypes.getf s T.fh) in
        (match N.Handle.get st.nodes fh with
         | Nodes.UnixHandle.File (fd, kind) ->
           let lwt_fd = Lwt_unix.of_unix_file_descr fd in
           with_fh_resource fh (function
             | Closed -> Lwt.fail Unix.(Unix_error (EBADF, "", ""))
             | Open ->
               setattr_fd valid s req lwt_fd kind
               >>= fun () ->
               Sys_stat_unix_lwt.fstat fd
               >>= fun stat ->
               get_ownership Xowner.fget fd
               >>= fun owner ->
               Lwt.return (Open, (stat, owner))
           )
           >>= fun stat_owner ->
           let attr_valid = UInt64.zero in
           let attr_valid_nsec = UInt32.zero in
           store_attr_of_stat_owner req stat_owner
           >>= fun store_attr ->
           IO.Out.write_reply req
             (Profuse.Out.Attr.create
                ~attr_valid ~attr_valid_nsec ~store_attr)
           >>= fun () ->
           Lwt.return st
         | Nodes.UnixHandle.Dir (_, _) ->
           (* This should never happen. *)
           Lwt.fail Unix.(Unix_error (EISDIR, "", ""))
        )
      else
        let { Nodes.data } = N.get st.nodes (nodeid req) in
        let path = Node.to_string data in
        let fd_flags = Fcntl.Oflags.(
          if valid.Valid.size
          then [O_WRONLY; O_SYMLINK]
          else [O_SYMLINK]
        ) in
        catch (fun () ->
          (* We'd like to atomically open an fd and then do our
             metadata operations on that handle so we don't race with
             other mutators. If this isn't supported, though, e.g. in
             the case of a socket file, we fall back to performing
             operations on the path. If we only have a single piece of
             metadata to update, we use the path functions directly.

             TODO: uid and gid update can be batched but aren't
          *)
          match Valid.to_string_list valid with
          | [] | [_] ->
            setattr_path valid s req path
          | _multi ->
            Fcntl_unix_lwt.open_ path fd_flags
            >>= fun fd ->
            let fd = Lwt_unix.of_unix_file_descr fd in
            finalize (fun () ->
              Lwt_unix.fstat fd
              >>= fun { Unix.st_kind } ->
              let kind = Sys_stat_unix.File_kind.of_unix st_kind in
              setattr_fd valid s req fd kind
            ) (fun () -> Lwt_unix.close fd)
        ) (function
          | Errno.Error { Errno.label = p; errno }
            when p = path
              && Errno.(List.mem EOPNOTSUPP errno
                        || List.mem ENOTSUP errno) ->
            setattr_path valid s req path
          | exn -> fail exn
        )
        >>= fun () ->
        getattr req st
    )

  end
end

let event_filter_pid = ref Unsigned.UInt32.zero

module EventFilter(F : Fuse_lwt.FS_LWT with type t = state)
  : Fuse_lwt.FS_LWT with type t = F.t = struct
  type t = F.t

  let string_of_state = F.string_of_state
  let string_of_nodeid = F.string_of_nodeid

  let log_error = log_error

  module Calls(IO : Fuse_lwt.IO_LWT)
    : Fuse_lwt.FS_IO_LWT with type t = t and module IO = IO = struct
    module Calls = F.Calls(IO)
    include Calls

    module Support = Fuse.Support(IO)

    let nodeid req = UInt64.to_int64 (Support.nodeid req)

    let current_syscall = ref None

    let set_syscall_segments segments = match !current_syscall with
      | None ->
        Log.err (fun f -> f "Syscall missing, couldn't set segments")
      | Some (_, syscall, path) ->
        current_syscall := Some (Some segments, syscall, path)

    let rec get_current_syscall () = match !current_syscall with
      | Some syscall -> Lwt.return syscall
      | None ->
        let open Lwt.Infix in
        Lwt_stream.get Volume.syscall_stream
        >>= function
        | None ->
          Log.err (fun f -> f "Couldn't get syscall for placation");
          Lwt.return (None, Volume.Truncate, "")
        | Some (path, syscall) ->
          let segs = Stringext.split ~on:'/' path in
          current_syscall := Some (Some segs, syscall, path);
          get_current_syscall ()

    let rec kind = Volume.(function
      | (Some (_ :: ((_ :: _) as segs)), syscall, path) ->
        current_syscall := Some (Some segs, syscall, path);
        Sys_stat.File_kind.DIR
      | (Some ([] | [_]), (Rmdir | Mkdir | Chmod_dir), _path) ->
        Sys_stat.File_kind.DIR
      | (Some ([] | [_]), Symlink, _path) ->
        Sys_stat.File_kind.LNK
      | (Some ([] | [_]), (Mknod_reg | Unlink | Truncate | Chmod), _path) ->
        Sys_stat.File_kind.REG
      | (Some _, Ping, _path) ->
        Log.err (fun f -> f "unexpected Ping in kind match");
        Sys_stat.File_kind.REG
      | (None, _, _) -> Sys_stat.File_kind.REG
    )

    (* This store_attr is only for placation responses. *)
    let store_attr kind req =
      let phost =
        Profuse.(req.chan.host.Host.sys_stat.Sys_stat.Host.file_kind)
      in
      let uid = Ctypes.getf req.Profuse.hdr Profuse.In.Hdr.T.uid in
      let gid = Ctypes.getf req.Profuse.hdr Profuse.In.Hdr.T.gid in
      let kind_code = Sys_stat.File_kind.to_code ~host:phost kind in
      Profuse.Struct.Attr.store
        ~ino:(UInt64.of_int 0) (* TODO: a future problem with FUSE ino flags *)
        ~size:UInt64.max_int
        ~blocks:UInt64.zero
        ~blksize:(UInt32.of_int 512)
        ~atime:UInt64.zero
        ~atimensec:UInt32.zero
        ~mtime:UInt64.zero
        ~mtimensec:UInt32.zero
        ~ctime:UInt64.zero
        ~ctimensec:UInt32.zero
        ~mode:(UInt32.of_int (kind_code + 0o777))
        ~nlink:UInt32.one
        ~uid
        ~gid
        ~rdev:UInt32.zero

    let rec drop_until el = function
      | Some ((next::_) as list) when next = el -> Some list
      | Some [] -> Some []
      | Some [_] | None -> None
      | Some (_::rest) -> drop_until el (Some rest)

    let return_error errno req st =
      let open Lwt.Infix in
      IO.Out.write_error log_error req errno
      >>= fun () -> Lwt.return st

    let enotrecoverable = return_error Errno.ENOTRECOVERABLE
    let enoent = return_error Errno.ENOENT

    let drop_event reason syscall path req st =
      Log.warn (fun f -> f "dropped event syscall %s for %s (%s)"
        (Volume.string_of_syscall syscall) path reason);
      current_syscall := None;
      enotrecoverable req st

    let ack req st =
      let open Lwt.Infix in
      IO.Out.write_ack req
      >>= fun () -> Lwt.return st

    (* This lookup is only for placation responses. *)
    let lookup name req st =
      let open Volume in
      let open Lwt.Infix in
      let { Nodes.data } = N.get st.nodes (nodeid req) in
      let target = Filename.concat (Node.to_string data) name in
      get_current_syscall ()
      >>= function
      | (Some (_::_), (Mkdir | Symlink | Mknod_reg), path)
        when path = target ->
        enoent req st (* The earth was formless and void *)
      | (Some [bad_seg],
         (Mkdir | Symlink | Mknod_reg as syscall),
         path) ->
        let reason = "bad create target "^bad_seg^" <> "^name in
        drop_event reason syscall path req st
      | (segs, syscall, path) ->
        match drop_until name segs with
        | None -> drop_event "bad lookup path" syscall path req st
        | segs_opt ->
          let parent = N.get st.nodes (nodeid req) in
          let node = N.lookup parent name in
          let nodeid = Unsigned.UInt64.of_int64 node.Nodes.id in
          let generation = Unsigned.UInt64.of_int64 node.Nodes.gen in
          let entry_valid = Unsigned.UInt64.zero in
          let entry_valid_nsec = Unsigned.UInt32.zero in
          let attr_valid = Unsigned.UInt64.zero in
          let attr_valid_nsec = Unsigned.UInt32.zero in
          let store_attr = store_attr (kind (segs_opt, syscall, path)) req in
          IO.Out.write_reply req
            (Profuse.Out.Entry.create ~nodeid ~generation
               ~entry_valid ~entry_valid_nsec
               ~attr_valid ~attr_valid_nsec
               ~store_attr)
          >>= fun () -> Lwt.return st

    (* This getattr is only for placation responses. *)
    let getattr req st =
      let attr_valid = UInt64.zero in
      let attr_valid_nsec = UInt32.zero in
      let open Lwt.Infix in
      get_current_syscall ()
      >>= fun syscall ->
      let store_attr = store_attr (kind syscall) req in
      IO.Out.write_reply req
        (Profuse.Out.Attr.create ~attr_valid ~attr_valid_nsec ~store_attr)
      >>= fun () -> Lwt.return st

    let complete_syscall parent segs syscall path req st expected ?name yes =
      let correct_syscall = Volume.(match expected with
        | Chmod | Chmod_dir -> syscall = Chmod || syscall = Chmod_dir
        | Rmdir | Unlink | Mkdir | Symlink | Mknod_reg | Truncate ->
          syscall = expected
        | Ping ->
          Log.err (fun f -> f "unexpected ping in syscall stream");
          false
      ) in
      if not correct_syscall
      then
        let reason = "mismatched syscall "^Volume.string_of_syscall expected in
        drop_event reason syscall path req st
      else
        let tpath = Node.to_string parent.Nodes.data in
        let tpath = match name with
          | Some name -> Filename.concat tpath name
          | None -> tpath
        in
        let open Lwt.Infix in
        if tpath = path
        then
          (set_syscall_segments [];
           yes req st
           >>= fun st ->
           current_syscall := None;
           enable_dcache parent name;
           Lwt.return st)
        else (drop_event ("bad target path "^tpath) syscall path req st)

    let ping_response syscall path req st =
      let open Lwt.Infix in
      Lwt_stream.peek ping
      >>= function
      | Some t0 ->
        let t1 = Mtime.(to_ns_uint64 (elapsed ())) in
        ping_reply_push (Some (t0, t1));
        current_syscall := None;
        enotrecoverable req st
      | None ->
        drop_event "ping FAILED" syscall path req st

    let event_response req st =
      let open Lwt.Infix in
      get_current_syscall ()
      >>= fun (segs, syscall, path) ->
      if path = "/host_docker_app/__ping__" && syscall = Volume.Symlink
      then ping_response syscall path req st
      else
        let parent = N.get st.nodes (nodeid req) in
        let complete = complete_syscall parent segs syscall path req st in
        Profuse.In.Message.(match req.Profuse.pkt with
          (* Placations *)
          | Rmdir name  -> complete Volume.Rmdir  ~name (fun req st ->
            N.unlink parent name;
            ack req st
          )
          | Unlink name -> complete Volume.Unlink ~name (fun req st ->
            N.unlink parent name;
            ack req st
          )
          | Mkdir (_,name) ->
            complete Volume.Mkdir ~name (lookup name)
          | Symlink (name,_) ->
            complete Volume.Symlink ~name (lookup name)
          | Setattr s ->
            let open Profuse.In.Setattr in
            let valid = Valid.of_uint32 (Ctypes.getf s T.valid) in
            if valid.Valid.size
            then complete Volume.Truncate getattr
            else if valid.Valid.mode
            then complete Volume.Chmod getattr
            else drop_event "unexpected setattr" syscall path req st
          | Mknod (c,name) ->
            complete Volume.Mknod_reg ~name (lookup name)

          (* Necessities *)
          | Interrupt _ ->
            (* The correct response for unimplemented interrupt support
               is to ignore *)
            Lwt.return st
          | Forget f ->
            forget (Ctypes.getf f Profuse.In.Forget.T.nlookup) req st
          | Lookup name -> lookup name req st
          | Destroy -> destroy req st

          (* Extra *)
          | message ->
            Log.warn (fun f -> f "UNEXPECTED event message %s"
                         (Profuse.In.Message.describe req));
            IO.Out.write_error log_error req Errno.ENOSYS
            >>= fun () -> Lwt.return st
        )

    let dispatch req t =
      let req_pid = Ctypes.getf req.Profuse.hdr Profuse.In.Hdr.T.pid in
      if Unsigned.UInt32.compare !event_filter_pid req_pid = 0
      then event_response req t
      else dispatch req t
  end
end

module HostFS = Fuse_lwt.Dispatch(FS)
module ServedFS = EventFilter(HostFS)

module FusePing = struct
  let session = ref 0
  let ping_file id =
    let log_dir = Osx_reporter.get_trace_dir () in
    Filename.concat log_dir (Printf.sprintf "osxfs.%04d.ping" id)

  let vsock_file id =
    let log_dir = Osx_reporter.get_trace_dir () in
    Filename.concat log_dir (Printf.sprintf "osxfs.%04d.vsock.ping" id)

  let write_event cstruct =
    let len = Cstruct.len cstruct in
    let buf = Bytes.create len in
    Cstruct.blit_to_bytes cstruct 0 buf 0 len;
    Lwt_unix.write !event_sock buf 0 len

  let ping () =
    let open Lwt.Infix in
    let now = Mtime.(to_ns_uint64 (elapsed ())) in
    ping_push (Some now);
    Volume.write_syscall "/host_docker_app/__ping__" write_event Volume.Symlink
    >>= fun () ->
    Lwt_stream.next ping_reply
    >>= fun (t0, event_t) ->
    if t0 = now
    then
      let event_rtt = Int64.sub event_t t0 in
      begin
        Lwt_stream.next ping_reply
        >>= fun (t0, error_t) ->
        if t0 = now
        then
          let error_rtt = Int64.sub error_t event_t in
          Lwt.return (event_rtt, error_rtt)
        else (
          Log.err
            (fun f -> f "Bad ping reply stream: wanted %Ld got %Ld" now t0);
          Lwt.return (0_L, 0_L)
        )
        end
    else (
      Log.err
        (fun f -> f "Bad ping reply stream: wanted %Ld got %Ld" now t0);
      Lwt.return (0_L, 0_L)
    )

  let vsock_ping () =
    let open Lwt.Infix in
    let now = Mtime.(to_ns_uint64 (elapsed ())) in
    ping_push (Some now);
    Volume.write_syscall "" write_event Volume.Ping
    >>= fun () ->
    Lwt_stream.next ping_reply
    >>= fun (t0, pong_t) ->
    if t0 = now
    then
      let rtt = Int64.sub pong_t t0 in
      Lwt.return rtt
    else (
      Log.err
        (fun f -> f "Bad ping reply stream: wanted %Ld got %Ld" now t0);
      Lwt.return 0_L
    )

  let rec session_files () =
    incr session;
    let filename = ping_file !session in
    let filename_v = vsock_file !session in
    let flags = Unix.([O_WRONLY; O_CREAT; O_EXCL]) in
    try
      (filename, Unix.(openfile filename flags 0o600)),
      (filename_v, Unix.(openfile filename_v flags 0o600))
    with Unix.Unix_error (Unix.EEXIST, "open", _) -> session_files ()

  let record () =
    let open Lwt.Infix in
    Lwt.async (fun () ->
      let rec reverb pings = function
        | 0 -> Lwt.return pings
        | k ->
          ping ()
          >>= fun result ->
          reverb (result::pings) (k - 1)
      in
      reverb [] 10000
      >>= fun pings ->
      (* TODO: lwtize *)
      let (filename, fd), (filename_v, fd_v) = session_files () in
      let oc = Unix.out_channel_of_descr fd in
      List.iter (fun (event_rtt, error_rtt) ->
        let line = Printf.sprintf "%Ld\t%Ld\n" event_rtt error_rtt in
        output_string oc line
      ) (List.rev pings);
      flush oc;
      Unix.close fd;

      let rec vsock_reverb pings = function
        | 0 -> Lwt.return pings
        | k ->
          vsock_ping ()
          >>= fun result ->
          vsock_reverb (result::pings) (k - 1)
      in
      vsock_reverb [] 10000
      >>= fun vpings ->
      let oc = Unix.out_channel_of_descr fd_v in
      List.iter (fun rtt ->
        let line = Printf.sprintf "%Ld\n" rtt in
        output_string oc line
      ) (List.rev vpings);
      flush oc;
      Unix.close fd_v;

      Log.warn
        (fun f -> f "ping results in %s; vsock results in %s"
            filename filename_v);
      Lwt.return_unit
    )
end

let rec dentry_of_path root = function
  | [] -> Some (UInt64.of_int64 root.Nodes.id, ".")
  | [name] -> Some (UInt64.of_int64 root.Nodes.id, name)
  | next::rest -> match Hashtbl.find root.Nodes.children next with
    | exception Not_found -> None
    | id -> dentry_of_path (N.get root.Nodes.space id) rest

module type OSXFS_SERVER =
  Fuse_lwt.SERVER_LWT
  with module IO = Fuse_lwt.IO
   and type t = ServedFS.t

let tracer = Trace.new_tracer ()

let mount server_module mounts =
  let module Server = (val server_module : OSXFS_SERVER) in
  let open Lwt.Infix in

  Lwt_list.map_s (fun mount ->
    let export_root = Mount.export_root mount in
    let state = make mount (match export_root with "/" -> "" | x -> x) in
    (* TODO: Add ability to trace capture the mount process *)
    (* TODO: Enable tracing for channels that are added when a session
       is active *)
    let args = [
      "allow_other";
      "max_read="^(string_of_int state.max_read);
      "subtype=osxfs";
    ] in
    let mnt = Mount.mount_path mount in
    Server.mount ~argv:[|"-o"; String.concat "," args|] ~mnt state
    >>= fun (req, state) ->
    Trace.add_channel tracer ~id:Profuse.(req.chan.id) mount;
    Lwt.return (req, state)
  ) mounts

let configure_syscall_push server_module clients =
  let module Server = (val server_module : OSXFS_SERVER) in
  let open Lwt.Infix in
  let syscall_push =
    let syscall_push = !Volume.syscall_push in
    function
    | None -> syscall_push None
    | Some (path, _syscall) as x ->
      let should_invalidate = is_path_dcached path in
      disable_dcache path;
      let after_invalidation =
        if should_invalidate
        then Lwt.join (List.map (fun (req, state) ->
          let root = N.root state.nodes in
          let root_path = Node.to_string root.Nodes.data in
          let root_path_len = String.length root_path in
          let path_len = String.length path in
          if path_len < root_path_len
          then Lwt.return_unit
          else if String.sub path 0 root_path_len <> root_path
          then Lwt.return_unit
          else
            let tail_len = path_len - root_path_len in
            let tail = String.sub path root_path_len tail_len in
            match Stringext.split ~on:'/' tail with
            | [] -> (* root *)
              Log.err (fun f -> f "pushing event for the root %s failed" path);
              Lwt.return_unit
            | ""::segments ->
              begin match dentry_of_path root segments with
                | None ->
                  (*Log.info (fun f -> f "no dentry for %s" path);*)
                  Lwt.return_unit
                | Some (parent, n) ->
                  let msg = Profuse.Out.Notify.Inval_entry.create parent n in
                  Server.IO.Out.write_notify req.Profuse.chan msg
                  >>= fun () ->
                  Server.IO.In.read_notify req.Profuse.chan ()
                  >>= function
                  | [] -> Lwt.return_unit
                  | errnos ->
                    if List.mem Errno.ENOENT errnos
                    then Lwt.return_unit
                    else
                      let errnos = List.map Errno.to_string errnos in
                      Log.err (fun f ->
                        f "While invalidating %s (%Ld / %s): %s"
                          path (UInt64.to_int64 parent) n
                          (String.concat ", " errnos)
                      );
                      Lwt.return_unit
              end
            | _ -> Lwt.return_unit
        ) clients)
        else Lwt.return_unit
      in
      after_invalidation
      >>= fun () ->
      syscall_push x
  in
  Volume.syscall_push := syscall_push

let configure_signals () =
  let is_tracing = ref false in
  Sys.(set_signal sigusr1 (Signal_handle (fun _ -> match !is_tracing with
    | true ->
      is_tracing := false;
      let dir = Trace.stop_tracing tracer in
      Log.info (fun f -> f "osxfs tracing: OFF %s" dir)
    | false ->
      is_tracing := true;
      let dir = Trace.start_tracing tracer in
      Log.info (fun f -> f "osxfs tracing: ON %s" dir)
  )));

  Sys.(set_signal sigusr2 (Signal_handle (fun _ -> FusePing.record ())));

  Sys.(set_signal sigint (Signal_handle (fun _ ->
    (* TODO: unmount *)
    Log.debug (fun f -> f "TODO: unmount");
    exit 1
  )));

  Sys.(set_signal sighup (Signal_handle (fun _ ->
    Log.info (fun f -> f "collecting and compacting");
    Gc.compact ()
  )))

module Volumes = Active_list.Make(Volume)
module VolumeServer = Server9p_unix.Make(LogInfo)(Volumes)

let transfuse_control_handler
    mnt_mod connect_path control_path volume_control_fs mounts =
  (module struct
    open Lwt.Infix

    module Mnt = (val mnt_mod : Fuse_lwt.MOUNT_LWT)
    module Server = Fuse_lwt.Server(Mnt)(ServedFS)(Fuse_lwt.IO)

    let mtab = Hashtbl.create 16

    let init ctl =
      (* Start listening for front-end commands *)
      Osx_socket.listen control_path
      >>= fun control_s ->
      Lwt.async (fun () -> Osxfs_control.serve_forever ctl mounts control_s);

      Transfuse.inotify_channel connect_path
      >>= fun (event_fd, actor_pid) ->
      Log.debug (fun f -> f "Event thread is %s"
                    (Unsigned.UInt32.to_string actor_pid));
      event_filter_pid := actor_pid;
      Volumes.set_context volume_control_fs Volume.({
        event_fd;
        mounts;
        ctl;
      });
      event_sock := event_fd;

      (* Create a mount table to track notify channels which come in
         asynchronously after the mount negotiation. *)
      List.iter (fun m ->
        Hashtbl.replace mtab (Mount.mount_path m) (Lwt_mvar.create_empty ())
      ) mounts;

      mount (module Server) mounts
      >>= fun clients ->

      (* Setup notify channels *)
      Lwt_list.iter_p (fun (req, state) ->
        let mount_point = Mount.mount_path state.mount in
        Lwt_mvar.take (Hashtbl.find mtab mount_point)
        >>= fun sock ->
        Log.info (fun f ->
          f "Negotiated transfuse notification channel for %s" mount_point
        );
        let nread () =
          let open Ctypes in
          let buf = allocate_n uint32_t ~count:1 in
          Socket.read_exactly sock (to_voidp buf) 4
          >>= fun () ->
          Lwt.return (!@ buf)
        in
        let nwrite_mutex = Lwt_mutex.create () in
        let nwrite p len =
          Lwt_mutex.with_lock nwrite_mutex (fun () ->
            Socket.write_exactly sock p len
          ) >>= fun () ->
          Lwt.return len
        in
        Fuse_lwt.set_socket Profuse.(req.chan.id) ~nread ~nwrite ();
        Lwt.return_unit
      ) clients
      >>= fun () ->

      configure_syscall_push (module Server) clients;
      configure_signals ();
      Control.continue ctl
      >>= fun () ->
      Lwt.join (List.map (fun (req, state) ->
        Server.serve_forever req.Profuse.chan state
        >>= fun _fs ->
        Lwt.return_unit
      ) clients)

    let log ctl level msg = match level, msg with
      | Control.Error,
        "Event symlink /host_docker_app/__ping__ error: State not recoverable\n" ->
        Lwt_stream.get ping
        >>= (function
          | Some t0 ->
            let t1 = Mtime.(to_ns_uint64 (elapsed ())) in
            ping_reply_push (Some (t0, t1));
            Lwt.return_unit
          | None ->
            Log.err (fun f -> f "ping stream empty!");
            Lwt.return_unit
        )
      | Control.Error, _ ->
        Log.err (fun f -> f "transfused: %s" (String.trim msg));
        Lwt.return_unit
      | Control.Notice, _ ->
        Log.info (fun f -> f "transfused: %s" (String.trim msg));
        Lwt.return_unit

    let pong ctl =
      Lwt_stream.get Volume.syscall_stream
      >>= (function
        | Some ("", Volume.Ping) -> Lwt.return_unit
        | Some _ | None ->
          Log.err (fun f -> f "syscall stream didn't have Ping!");
          Lwt.return_unit
      )
      >>= fun () ->
      Lwt_stream.get ping
      >>= (function
        | Some t0 ->
          let t1 = Mtime.(to_ns_uint64 (elapsed ())) in
          ping_reply_push (Some (t0, t1));
          Lwt.return_unit
        | None ->
          Log.err (fun f -> f "ping stream empty!");
          Lwt.return_unit
      )

    let set_notify_channel sock mount_point =
      Lwt_mvar.put (Hashtbl.find mtab mount_point) sock

  end : Control.HANDLER)

let serve connect_path address db control_path volume_control_path debug =
  Osx_reporter.install ~stdout:debug;

  let conf = (module struct
               let max_write = max_write
               let connect_path = connect_path
               let address = address
             end : Transfuse.CONF)
  in
  let module Conf = (val conf) in
  let module Mnt = Transfuse.Mount(Conf) in
  let open Lwt.Infix in

  let t =
    (* Load the export:mount map from the database *)
    Mount.load db
    >>= fun mounts ->

    (* Start volume approval server *)
    let volume_control_fs = Volumes.make () in
    Osx_socket.listen volume_control_path
    >>= fun volume_s ->

    (* Start listening for bind mount requests *)
    let volumeserver = VolumeServer.of_fd volume_control_fs volume_s in
    Lwt.async (fun () -> VolumeServer.serve_forever volumeserver);

    Transfuse.listen_forever address
      (transfuse_control_handler
         (module Mnt)
         connect_path
         control_path
         volume_control_fs
         mounts)

  in
  Lwt.async_exception_hook := begin fun exn ->
    let exn_str = Printexc.to_string exn in
    Log.err (fun f -> f "Fatal unexpected exception: %s" exn_str);
    exit 1
  end;
  try ignore (Lwt_main.run t); `Ok () with
  | Failure e ->
    Log.err (fun f -> f "Fatal unexpected failure: %s" e);
    `Error (false, e)
  | e         ->
    let exn_str = Printexc.to_string e in
    Log.err (fun f -> f "Fatal unexpected exception: %s" exn_str);
    `Error (false, exn_str)

open Cmdliner

let address =
  let doc =
    "Address of the fuse server: unix:path"
  in
  Arg.(value & opt string "unix:fuse.sock" & info [ "address"; "a" ] ~doc)

let connect =
  let doc =
    "Socket path of the callback proxy"
  in
  Arg.(value & opt string "connect.sock" & info [ "connect"; "c" ] ~doc)

let database =
  let doc = "Path to database socket" in
  Arg.(required & opt (some string) None & info [ "database" ] ~doc)

let control_path =
  let doc = "Path to osxfs control socket" in
  Arg.(required & opt (some string) None & info ["control"] ~doc)

let volume_control_path =
  let doc = "Path to volume control socket" in
  Arg.(required & opt (some string) None & info [ "volume-control" ] ~doc)

let control_path =
  let doc = "Path to osxfs control socket" in
  Arg.(required & opt (some string) None & info ["control"] ~doc)

let debug =
  let doc = "Verbose debug logging to stdout" in
  Arg.(value & flag & info [ "debug" ] ~doc)

let serve_cmd =
  let doc = "Serve an OS X directory over transfuse" in
  let man = [
    `S "DESCRIPTION";
    `P "Listen for transfuse connections and serve the named filesystem.";
  ] in
  Term.(ret(pure serve
            $ connect
            $ address
            $ database
            $ control_path
            $ volume_control_path
            $ debug)),
  Term.info "osxfs" ~doc ~man

;;
Sys.set_signal Sys.sigpipe Sys.Signal_ignore;
match Term.eval serve_cmd with
| `Error _ -> exit 1
| _ -> exit 0
