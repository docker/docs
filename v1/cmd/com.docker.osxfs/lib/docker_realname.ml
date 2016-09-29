(* What's the name of the host object you are trying to mount?
   Docker_realname.resolve will tell you according to docker+OS X semantics.

   Seriously.
*)

open Lwt.Infix

exception Error of string

let get_name_and_type path =
  let open Osx_attr in
  let name_type = Common.(Query.[Common NAME; Common OBJTYPE]) in
  Osx_attr_lwt.getlist ~no_follow:true name_type path
  >>= function
  | [ Value.Common (Common.NAME, name);
      Value.Common (Common.OBJTYPE, t) ] ->
    Lwt.return ((name : string), (t : Osx_attr.Vnode.Vtype.t))
  | [ Value.Common (Common.OBJTYPE, t);
      Value.Common (Common.NAME, name) ] ->
    Lwt.return ((name : string), (t : Osx_attr.Vnode.Vtype.t))
  | _ -> Lwt.fail (Error ("Could not get name and/or type for "^path))

let absolute_path_rel_to_root path =
  let len = String.length path in
  if len > 0 && path.[0] = '/'
  then String.sub path 1 (len - 1)
  else raise (Error ("path "^path^" is not absolute"))

let safe_tl = function
  | [] -> []
  | _::rest -> rest

let string_of_stack stack = "/"^(String.concat "/" (List.rev stack))

let resolve path = Lwt.catch (fun () ->
  let rec step realname = function
    | [] -> Lwt.return (string_of_stack realname)
    | (""|".")::rest -> step realname rest
    | ".."::rest -> step (safe_tl realname) rest
    | seg::rest ->
      let open Osx_attr.Vnode.Vtype in
      let obj = seg::realname in
      let obj_s = string_of_stack obj in
      Lwt.catch
        (fun () ->
           get_name_and_type obj_s
           >>= function
           | (name, VLNK) ->
             Lwt_unix.readlink obj_s
             >>= fun target ->
             if target.[0] = '/'
             then
               let path = absolute_path_rel_to_root target in
               step [] (Stringext.split path ~on:'/')
             else step realname ((Stringext.split target ~on:'/')@rest)
           | (name, _) -> step (name::realname) rest
        )
        (function
          | Errno.Error { Errno.errno } when List.mem Errno.ENOENT errno ->
            Lwt.return (string_of_stack (List.rev_append rest obj))
          | exn -> Lwt.fail exn
        )
  in
  let path = absolute_path_rel_to_root path in
  step [] (Stringext.split path ~on:'/')
) (function
  | Error message -> Lwt.fail (Error message)
  | exc -> Lwt.fail (Error (Printexc.to_string exc))
)
