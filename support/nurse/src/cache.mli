

type 'a error = ('a, [ `Msg of string ]) Result.result

type diagnostic_id = string

type timestamp = string

type path = string

val get: diagnostic_id -> timestamp -> path error Lwt.t
(** Return a local path of a decompressed archive *)
