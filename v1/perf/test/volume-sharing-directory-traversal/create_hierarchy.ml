(* Create a directory/file hierarchy from a sexp *)

let parse lexbuf =
  Sexplib.Parser.sexp
    (fun lexbuf -> Sexplib.Lexer.main lexbuf)
    lexbuf

let parse_string s = parse (Lexing.from_string s)
let parse_file filename = parse (Lexing.from_channel (open_in filename))

(*
   (A (B C D) E)
~>
   A/
    B/
      C
      D
    E
*)
let rec create ~parent sexp =
  let open Sexplib.Type in
  match sexp with
  | Atom name ->
    Unix.close @@
    Unix.openfile (Filename.concat parent name)
      Unix.[O_CREAT]
      0o755
  | List [] ->
    ()
  | List (Atom dir :: contents) ->
    begin
      let full_path = Filename.concat parent dir in
      Unix.mkdir full_path 0o755;
      List.iter (create ~parent:full_path) contents
    end
  | List (List _ :: _) ->
    Printf.ksprintf failwith "Missing directory name under %s" parent

let main ~parent filename =
  create ~parent (parse_file filename)

let () = main ~parent:"." Sys.argv.(1)
