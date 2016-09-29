
val run: string -> string list -> (string, [ `Exit of int | `Msg of string ]) Result.result Lwt.t
(** [run path args] runs a process and returns stdout *)
