
type t = {
  in_file: string;                       (** Name of file within the diagnostic archive *)
  regexp: string;                        (** If this regexp matches, the problem probably exists *)
  label: string;                         (** A short label (e.g. "qcow2-corruption") *)
  description: string;                   (** Description suitable for posting as a github comment *)
  link_to_issues: (string * int) list;   (** (repo, issue) we should link this to *)
  link_from_issues: (string * int) list; (** (repo, issue) we should link to this *)
} [@@deriving sexp]
(** A well-known problem *)

val compiled_in: t list
(** Problems compiled into this binary (as opposed to downloaded) *)
