type t

val make : export_path:string -> mount_path:string -> t Lwt.t

val export_path : t -> string
val export_root : t -> string
val mount_path : t -> string

val segment_name : t -> string

val load : string -> t list Lwt.t
