
type t
(** A diagnostic archive *)

val get_diagnostic_id: t -> string
(** The unique ID which represents the software installation *)

val get_timestamp: t -> string
(** The time the diagnostics were generated *)

type 'a error = ('a, [ `Msg of string]) Result.result

val openarchive: string -> t error Lwt.t
(** [openarchive filename] opens the archive and prepares it for queries. *)

val close: t -> unit Lwt.t
(** [close t] closes the archive, releasing associated resources. *)

module Symptom: sig
  type symptom = {
    problem: Problem.t;    (** the problem we found *)
    archive: t;            (** the archive we found it in *)
    contexts: string list; (** the log file context(s) *)
  }
  (** Symptom(s) of a specific problem in a specific archive *)

  val to_markdown: symptom -> string
  (** Render a report as markdown *)
end

val analyse: t -> Symptom.symptom list Lwt.t
(** [analyse t] looks for symptoms of known problems in the archive. *)
