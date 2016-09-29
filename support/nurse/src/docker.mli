

val inspect: string -> (string, [ `Exit of int | `Msg of string ]) Result.result Lwt.t
(** Same as the output of `docker inspect` *)

val run: string list -> (string, [ `Exit of int | `Msg of string ]) Result.result Lwt.t
(** Same as the output of `docker run <args>` *)
