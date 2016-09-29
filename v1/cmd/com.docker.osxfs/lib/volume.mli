
type t
(** The client requests one of these *)

type syscall =
  | Ping
  | Rmdir
  | Unlink
  | Mkdir
  | Symlink
  | Truncate
  | Chmod_dir
  | Chmod
  | Mknod_reg

val write_syscall : string -> (Cstruct.t -> int Lwt.t) -> syscall -> unit Lwt.t

val syscall_push : ((string * syscall) option -> unit Lwt.t) ref

val syscall_stream : (string * syscall) Lwt_stream.t

val string_of_syscall : syscall -> string

val to_string: t -> string
val of_string: string -> (t, [ `Msg of string ]) Result.result

type context = {
  event_fd : Lwt_unix.file_descr;
  mounts : Mount.t list;
  ctl : Control.t;
}
(** The context in which a [t] is [start]ed, for example a TCP/IP stack *)

val start: context Active_list.Var.t -> t -> (t, [ `Msg of string ]) Result.result Lwt.t

val stop: t -> unit Lwt.t

type key
(** Some unique primary key *)

module Map: Map.S with type key = key

val get_key: t -> key

val description_of_format: string
