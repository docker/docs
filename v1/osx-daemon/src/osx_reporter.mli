val install: stdout:bool -> unit
(** [install stdout] installs a log reporter. If stdout is true, then logs
    are sent only to stdout, otherwise logs are sent to ASL *)

val get_trace_dir: unit -> string
(** [get_trace_dir ()] returns the directory where this process should write traces.
    The directory will be created by this function if it does not already
    exist. *)
