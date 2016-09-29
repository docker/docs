open Lwt.Infix
open Astring

module KV_RO = Tar_mirage.Make_KV_RO(Block)

type t = {
  block: Block.t;
  kv_ro: KV_RO.t;
  diagnostic_id: string;
  timestamp: string;
}

module Symptom = struct
  type symptom = {
    problem: Problem.t;
    archive: t;
    contexts: string list;
  }

  let to_markdown { problem; archive; contexts } =
    let quote s =
      let lines = String.cuts ~sep:"\n" s in
      String.concat ~sep:"\n" @@ List.map (fun x -> "> " ^ x) lines in
    Printf.sprintf "Detected symptom of problem '%s' in %s/%s.\n\n%s\n\nMay be related to %s\n\nThe following log matches:\n %s"
      problem.Problem.label archive.diagnostic_id archive.timestamp
      problem.Problem.description
      (String.concat ~sep:" " @@ List.map (fun (repo, no) -> repo ^ "#" ^ (string_of_int no)) problem.Problem.link_to_issues)
      (String.concat ~sep:"\nand\n" @@ List.map quote contexts)
end

let get_diagnostic_id t = t.diagnostic_id
let get_timestamp t = t.timestamp

type 'a error = ('a, [ `Msg of string]) Result.result

let openarchive filename =
  (* Add a `list` function to `Tar_mirage.Make_KV_RO`? *)
  Lwt_unix.openfile filename [ Lwt_unix.O_RDONLY ] 0
  >>= fun fd ->
  Lwt.finalize
    (fun () ->
       Tar_lwt_unix.Archive.list fd
    ) (fun () ->
        Lwt_unix.close fd
      )
  >>= fun headers ->

  Block.connect filename
  >>= function
  | `Error e ->
    Lwt.return (Result.Error (`Msg (Mirage_block.Error.string_of_error e)))
  | `Ok block ->
    Lwt.catch
      (fun () ->
         KV_RO.connect block
         >>= function
         | `Ok kv_ro ->
           match headers with
           | [] ->
             Lwt.return (Result.Error (`Msg (Printf.sprintf "Archive %s has no files" filename)))
           | first :: _rest ->
             begin match String.cuts ~sep:"/" first.Tar_lwt_unix.Header.file_name with
               | diagnostic_id :: timestamp :: _ ->
                 Lwt.return (Result.Ok { block; kv_ro; diagnostic_id; timestamp })
               | _ ->
                 Lwt.return (Result.Error (`Msg (Printf.sprintf "Archive %s file %s does not obey diagnostic_id/timestamp convention" filename first.Tar_lwt_unix.Header.file_name)))
             end
      ) (fun e ->
          Lwt.return (Result.Error (`Msg (Printf.sprintf "Unexpected exception %s while processing %s" (Printexc.to_string e) filename)))
        )
    >>= function
    | Result.Error x ->
      Block.disconnect block
      >>= fun () ->
      Lwt.return (Result.Error x)
    | Result.Ok x ->
      Lwt.return (Result.Ok x)

let close t =
  Block.disconnect t.block

let analyse t =
  Lwt_list.map_s
    (fun problem ->
      let key = t.diagnostic_id ^ "/" ^ t.timestamp ^ "/" ^ problem.Problem.in_file in
      KV_RO.size t.kv_ro key
      >>= function
      | `Error _ ->
        Logs.err (fun f -> f "Failed to find %s" key);
        Lwt.return []
      | `Ok size ->
        KV_RO.read t.kv_ro key 0 (Int64.to_int size)
        >>= function
        | `Error _ -> assert false
        | `Ok buffers ->
          let s = String.concat ~sep:"" (List.map Cstruct.to_string buffers) in
          (* print_endline s; *)
          try
            let re = Re_perl.(compile @@ re ~opts:[`Multiline] problem.Problem.regexp) in
            let groups = Re.exec re s in
            (* At least one instance of the problem detected. Compute a context
               substring for each one and then merge together (to avoid redundancy) *)
            let rec loop i = match Re.Group.start groups i, Re.Group.stop groups i with
              | start, stop -> (start, stop) :: (loop (i + 1))
              | exception Not_found -> [] in
            let idxs = List.rev @@ loop 0 in
            (* Aim to include at least [context_size] characters surrounding the match *)
            let context_size = 1000 in
            (* Snap to newlines and the edges of the string *)
            let rec snap_to_newline delta from =
              if from < 0 then 0
              else if from >= (String.length s) then String.length s - 1
              else if s.[from] = '\n' then from
              else snap_to_newline delta (from + delta) in
            let contexts = List.map (fun (start, stop) ->
              snap_to_newline (-1) (start - context_size / 2),
              snap_to_newline (+1) (stop + context_size / 2)
            ) idxs in
            (* Merge overlapping contexts to avoid redundant printing. *)
            let merge = function
              | [] -> []
              | at_least_one ->
                let sorted = List.sort (fun a b -> compare (fst a) (fst b)) at_least_one in
                let rec loop (a_start, a_end) = function
                  | [] -> [ a_start, a_end ]
                  | (b_start, b_end) :: next when b_start <= a_end -> loop (a_start, b_end) next
                  | b :: next -> (a_start, a_end) :: (loop b next) in
                loop (List.hd sorted) (List.tl sorted) in
            let contexts = merge contexts in
            let contexts = List.map
              (fun (context_start, context_end) ->
                let context_length = context_end - context_start in
                String.with_range ~first:context_start ~len:context_length s
              ) contexts in
            Lwt.return [ { Symptom.problem; archive = t; contexts } ]
          with
          | Not_found -> Lwt.return [] (* problem not detected *)
          | Re_perl.Parse_error -> failwith (Printf.sprintf "Failed to parse perl-style regexp: %s" problem.Problem.regexp)
    ) Problem.compiled_in
    >>= fun x ->
    Lwt.return (List.concat x)
