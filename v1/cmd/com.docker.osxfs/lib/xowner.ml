
open Lwt.Infix

let xattr_key = "com.docker.owner"

module Version_0 = struct

  type t = {
    uid : int option;
    gid : int option;
  }

  let vid = 0

  let to_string { uid; gid } = match uid, gid with
    | None,     None     -> ":"
    | Some uid, None     -> Printf.sprintf "%d:"   uid
    | None,     Some gid -> Printf.sprintf ":%d"   gid
    | Some uid, Some gid -> Printf.sprintf "%d:%d" uid gid

  let of_string s = match Stringext.split ~max:2 s ~on:':' with
    | [] | [_] | _::_::_::_ -> None
    | [""; ""] -> Some { uid = None; gid = None }
    | [uid; ""] ->
      let uid = try Some (int_of_string uid) with _ -> None in
      Some { uid; gid = None }
    | [""; gid] ->
      let gid = try Some (int_of_string gid) with _ -> None in
      Some { uid = None; gid }
    | [uid; gid] ->
      let uid = try Some (int_of_string uid) with _ -> None in
      let gid = try Some (int_of_string gid) with _ -> None in
      Some { uid; gid }

  let overlay bottom top = {
    uid = (match top.uid with None -> bottom.uid | Some _ -> top.uid);
    gid = (match top.gid with None -> bottom.gid | Some _ -> top.gid);
  }

  let empty = { uid = None; gid = None }
end

type t = Version_0.t = {
  uid : int option;
  gid : int option;
}

type 'a generic_get =
  ?show_compression:bool -> ?size:int -> 'a -> string -> string option Lwt.t

type 'a generic_set =
  ?create:bool -> ?replace:bool -> 'a -> string -> string -> unit Lwt.t

let to_string owner =
  Printf.sprintf "%d:%s" Version_0.vid (Version_0.to_string owner)

let of_string attr = match Stringext.split ~max:2 attr ~on:':' with
  | [] | [_] | _::_::_::_ -> None
  | [version; rest] ->
    if version = string_of_int Version_0.vid
    then Version_0.of_string rest
    else None

let overlay = Version_0.overlay

let empty = Version_0.empty

let set_generic (get : 'a generic_get) (set : 'a generic_set) ?uid ?gid file =
  let owner = { uid; gid } in
  begin match owner with
    | { uid = None; gid = None } -> Lwt.return ""
    | { uid = Some _; gid = Some _ } -> Lwt.return (to_string owner)
    | { uid = Some _; gid = None } | { uid = None; gid = Some _ } ->
      get ~size:16 file xattr_key
      >>= function
      | None -> Lwt.return (to_string owner)
      | Some xowner ->
        let bottom = match of_string xowner with None -> empty | Some x -> x in
        Lwt.return (to_string (overlay bottom owner))
  end >>= fun value ->
  set file xattr_key value

let set  =
  set_generic
    (Osx_xattr_lwt.get ~no_follow:true)
    (Osx_xattr_lwt.set ~no_follow:true)
let fset =
  set_generic Osx_xattr_lwt.fget Osx_xattr_lwt.fset

let get_generic (get : 'a generic_get) file =
  get ~size:16 file xattr_key
  >>= function
  | None -> Lwt.return empty
  | Some xattr -> match of_string xattr with
    | None -> Lwt.return empty
    | Some owner -> Lwt.return owner

let get  = get_generic (Osx_xattr_lwt.get ~no_follow:true)
let fget = get_generic Osx_xattr_lwt.fget
