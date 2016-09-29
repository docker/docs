
module StringMap = Map.Make(String)

module Fs = struct

  type stream = {
    leftover : Cstruct.t option;
    reader   : Cstruct.t Lwt_stream.t;
    push     : Cstruct.t option -> unit;
  }

  type node =
    | Dir of tree StringMap.t
    | Stream of stream
    | File of Cstruct.t
  and tree = {
    name : string;
    qid : Protocol_9p.Types.Qid.t;
    writable : bool;
    removable : bool;
    mutable parent : tree option;
    mutable contents : node;
  }

  type t = {
    root         : tree;
    connections  : tree;
    events       : tree;
    event_stream : stream;
  }

  let qid_path = ref 0_L

  let next_qid flags =
    let id = !qid_path in
    qid_path := Int64.(add one !qid_path);
    Protocol_9p.Types.Qid.({ flags; version = 0_l; id; })

  let new_stream ?(removable=false) ?(writable=false) name =
    let reader, push = Lwt_stream.create () in
    let stream = {
      leftover = None;
      reader;
      push;
    } in
    {
      name;
      parent = None;
      qid = next_qid [];
      writable;
      removable;
      contents = Stream stream;
    }, stream

  let new_dir ?(children=[]) name =
    let dir = {
      name;
      parent = None;
      qid = next_qid Protocol_9p.Types.Qid.([ Directory ]);
      writable = false;
      removable = false;
      contents = Dir StringMap.empty;
    } in
    let children = List.fold_left (fun map tree ->
      tree.parent <- Some dir;
      StringMap.add tree.name tree map
    ) StringMap.empty children in
    dir.contents <- Dir children;
    dir

  let new_file name contents = {
    name;
    parent = None;
    qid = next_qid [];
    writable = false;
    removable = false;
    contents = File contents;
  }

  let insert_child dir node = match dir.contents with
    | Dir children ->
      let children = StringMap.add node.name node children in
      node.parent <- Some dir;
      dir.contents <- Dir children;
    | _ -> assert false (* TODO: :-( *)

  let readme = Cstruct.of_string {|
Active connections directory
----------------------------

Every accepted connection is represented by a numbered subdirectory.
Within each subdirectory there are 2 files:

 - read: reading this file reads from the socket
 - write: writing to this file writes to the socket

To close the connection, remove the \"write\" file.

Note: it is not very sensible for two processes to share the same
connection -- you have been warned!
|}

  let create () =
    let readme = new_file "README" readme in
    let events, event_stream = new_stream "events" in
    let connections = new_dir "connections" in
    let root = new_dir ~children:[events; connections; readme] "/" in
    { root; connections; events; event_stream }

end

type t = {
  read  : int -> Cstruct.t Lwt.t;
  write : Cstruct.t -> int Lwt.t;
}

let host = Profuse.Host.linux_4_0_5

let rec read_stream tree s count =
  let open Fs in
  match s.leftover with
  | None ->
    let open Lwt in
    Lwt_stream.get s.reader
    >>= fun leftover ->
    read_stream tree { s with leftover } count
  | Some ({ Cstruct.len = 0 } as leftover) ->
    tree.contents <- Stream { s with leftover = None };
    Lwt.return leftover
  | Some ({ Cstruct.len } as leftover) when len > count ->
    let data, leftover = Cstruct.split leftover count in
    tree.contents <- Stream { s with leftover = Some leftover };
    Lwt.return data
  | Some data ->
    let leftover = Some (Cstruct.create 0) in
    tree.contents <- Stream { s with leftover };
    Lwt.return data

let conn_id = ref 1
let open_socket protocol fs =
  let dir_name = string_of_int !conn_id in
  let read_file, rstream = Fs.new_stream "read" in
  (*read.reader.onclose = func() {
    s.connections.removeChild(dirName)
    }*)
  let conn_dir = Fs.new_dir dir_name in
  Fs.insert_child conn_dir read_file;
  let write_file, wstream =
    Fs.new_stream ~removable:true ~writable:true "write"
  in
  Fs.insert_child conn_dir write_file;
  Fs.insert_child Fs.(fs.connections) conn_dir;
  let announcement = Printf.sprintf "%s%s\n" protocol dir_name in
  Fs.(fs.event_stream.push) (Some (Cstruct.of_string announcement));

  let read = read_stream write_file wstream in
  let write buf =
    rstream.Fs.push (Some buf);
    Lwt.return (Cstruct.len buf)
  in
  incr conn_id;
  { read; write; }

(* TODO: fs should be independent between connections *)
let fs = ref (Fs.create ())

module Server : sig
  include Protocol_9p.Filesystem.S

  val make: unit -> t
end = struct
  open Protocol_9p

  type t = unit

  let make () = ()

  type connection = {
    info: Protocol_9p.Info.t;
    fids: Fs.tree Types.Fid.Map.t ref;
  }

  let connect () info =
    let fids = ref Types.Fid.Map.empty in
    { fids; info }

  let disconnect _ _info =
    Lwt.return ()

  module Error = struct
    let badfid = Lwt.return (Response.error "fid not found")
    let badwalk = Lwt.return (Response.error "bad walk") (* TODO: ? *)

    let enoent = Lwt.return (Response.error "file not found")
    let eperm  = Lwt.return (Response.error "permission denied")
  end

  let return x = Lwt.return (Result.Ok x)

  let attach connection ~cancel { Request.Attach.fid } =
    let socket_fs = Fs.create () in
    fs := socket_fs;
    connection.fids := Types.Fid.Map.add fid socket_fs.Fs.root !(connection.fids);
    return { Response.Attach.qid = Fs.(socket_fs.root.qid) }

  let walk connection ~cancel { Request.Walk.fid; newfid; wnames } =
    try
      let from = Types.Fid.Map.find fid !(connection.fids) in
      let from, wqids = List.fold_left (fun (from,qids) -> function
        | "" | "." -> from, from.Fs.qid::qids
        | ".."     -> begin match from.Fs.parent with
          | Some parent -> parent, parent.Fs.qid::qids
          | None -> failwith "ENOENT"
        end
        | name     -> Fs.(match from.contents with
          | Stream _ | File _ -> failwith "BADWALK" (* TODO: FIXME? *)
          | Dir children ->
            try
              let child = StringMap.find name children in
              child, child.qid::qids
            with Not_found -> failwith "ENOENT"
        )
      ) (from,[]) wnames in
      connection.fids := Types.Fid.Map.add newfid from !(connection.fids);
      let wqids = List.rev wqids in
      return { Response.Walk.wqids }
    with
    | Not_found -> Error.badfid
    | Failure "ENOENT" -> Error.enoent
    | Failure "BADWALK" -> Error.badwalk

  let clunk _connection ~cancel _ =
    Lwt.return (Result.Ok ())

  let open_ connection ~cancel { Request.Open.fid; mode } =
    try
      let tree = Types.Fid.Map.find fid !(connection.fids) in
      let qid = tree.Fs.qid in
      let iounit = 32768_l in
      return { Response.Open.qid; iounit }
    with Not_found -> Error.badfid

  let stat_tree tree =
    let open Fs in
    let is_directory = match tree.contents with Dir _ -> true | _ -> false in
    let exec = if is_directory then [ `Execute ] else [] in
    let perms =
      `Read :: (if tree.writable then [ `Write ] else [] ) @ exec
    in
    Types.({
      Stat.ty = 0xFFFF;
      dev     = Int32.(neg one);
      qid     = tree.qid;
      mode    = FileMode.make
          ~owner:perms ~group:perms ~other:perms ~is_directory ();
      atime   = 1146711721l;
      mtime   = 1146711721l;
      length  = 0_L; (* TODO: wrong for regular files *)
      name    = tree.name;
      uid     = "uid";
      gid     = "gid";
      muid    = "muid";
      u       = None;
    })

  let errors_to_client = Result.(function
    | Error (`Msg msg) -> Error { Response.Err.ename = msg; errno = None }
    | Ok _ as ok -> ok
  )

  let read connection ~cancel { Request.Read.fid; offset; count } =
    let count = Int32.to_int count in
    try
      let tree = Types.Fid.Map.find fid !(connection.fids) in
      Fs.(match tree.contents with
        | File data ->
          let offset = Int64.to_int offset in
          let len = min count Cstruct.(len data - offset) in
          let data = Cstruct.sub data offset len in
          return { Response.Read.data }
        | Stream s ->
          let open Lwt in
          read_stream tree s count
          >|= fun data ->
          Result.Ok { Response.Read.data }
        | Dir children ->
          let open Infix in
          let rec write off rest = function
            | [] -> return off
            | x :: xs ->
              let stat = stat_tree x in
              let n = Types.Stat.sizeof stat in
              if off < offset
              then write Int64.(add off (of_int n)) rest xs
              else if Cstruct.len rest < n then return off
              else
                Lwt.return (Types.Stat.write stat rest)
                >>*= fun rest ->
                write Int64.(add off (of_int n)) rest xs
          in
          let buf = Cstruct.create count in
          let t =
            write 0_L buf (List.rev_map snd (StringMap.bindings children))
            >>*= fun offset' ->
            let data =
              Cstruct.sub buf 0 Int64.(to_int (max 0_L (sub offset' offset)))
            in
            return { Response.Read.data }
          in
          Lwt.(t >>= fun x -> return (errors_to_client x))
      )
    with Not_found -> Error.badfid

  let stat connection ~cancel { Request.Stat.fid } =
    try
      let tree = Types.Fid.Map.find fid !(connection.fids) in
      return { Response.Stat.stat = stat_tree tree }
    with Not_found -> Error.badfid

  let create _connection ~cancel _ = Error.eperm

  let write connection ~cancel { Request.Write.fid; offset; data } =
    let open Fs in
    try
      let tree = Types.Fid.Map.find fid !(connection.fids) in
      if not tree.writable
      then Error.eperm
      else match tree.contents with
        | Dir _ -> Error.eperm
        | File file ->
          let count = Cstruct.len data in
          Cstruct.blit data 0 file (Int64.to_int offset) count;
          return { Response.Write.count = Int32.of_int count }
        | Stream { Fs.push } ->
          push (Some data);
          return { Response.Write.count = Int32.of_int (Cstruct.len data) }
    with Not_found -> Error.badfid

  let remove connection ~cancel { Request.Remove.fid } =
    let open Fs in
    try
      let tree = Types.Fid.Map.find fid !(connection.fids) in
      if not tree.writable
      then Error.eperm
      else match tree.parent with
        | None -> Error.eperm
        | Some p -> match p.contents with
          | File _ | Stream _ -> Error.eperm
          | Dir children ->
            p.contents <- Dir (StringMap.remove tree.name children);
            return ()
    with Not_found -> Error.badfid

  let wstat _connection ~cancel _ = Error.eperm
end

let inotify_to_socket9p () =
  let open Lwt.Infix in
  let socket9p = open_socket "e" !fs in
  socket9p.read 16
  >>= fun cstruct ->
  let actor_pid =
    Unsigned.UInt32.of_string (String.trim (Cstruct.to_string cstruct))
  in
  Lwt.return (socket9p, actor_pid)

let mount_to_socket9p ~mnt ~argv =
  let open Lwt in
  let opts = String.concat "\000" (Array.to_list argv @ [mnt]) in
  let sz = String.length opts in
  let socket9p = open_socket "m" !fs in
  socket9p.write (Cstruct.of_string opts)
  >>= fun n ->
  (if sz <> n then failwith "Couldn't write FUSE mount options");
  return socket9p

module type CONF = sig
  val max_write : int
end

module Mount(Conf : CONF)(FS:Fuse_lwt.FS_LWT)(IO: Fuse_lwt.IO_LWT)
  : Fuse.MOUNT_IO
    with module IO = IO
     and type t = FS.t
= struct
  module IO = IO
  type t = FS.t

  let mount negotiate ~argv ~mnt state =
    let open IO in
    mount_to_socket9p ~mnt ~argv
    >>= fun socket9p ->
    let socket = Fuse_lwt.new_socket
        ~read:(fun sz ->
          socket9p.read sz
          >>= fun cstruct ->
          let ba = Cstruct.to_bigarray cstruct in
          let ptr = Ctypes.(bigarray_start array1 ba) in
          let ptr = Ctypes.(coerce (ptr char) (ptr uint8_t)) ptr in
          return (Ctypes.CArray.from_ptr ptr (Cstruct.len cstruct))
        )
        ~write:(fun p sz ->
          let p = Ctypes.(coerce (ptr uint8_t) (ptr char)) p in
          let bigarray =
            Ctypes.(bigarray_of_ptr array1) sz Bigarray.Char p
          in
          socket9p.write (Cstruct.of_bigarray bigarray)
        )
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
        let opcode = Opcode.to_string (Ctypes.getf hdr Hdr.T.opcode) in
        let msg =
          Printf.sprintf "Unexpected opcode %s <> FUSE_INIT" opcode
        in
        IO.fail (ProtocolError (init_fs, msg))
    ))

end
