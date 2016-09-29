#! /usr/bin/env ocamlscript
Ocaml.packs := [ "unix"; "stringext" ];
Ocaml.sources := [ "../../lib/functions.ml"]
--
open Functions

let (/) = Filename.concat

module Normal = struct
  type t = {
    sigma_x: float;
    sigma_xx: float;
    n: int;
  }
  let empty = { sigma_x = 0.; sigma_xx = 0.; n = 0 }
  exception Insufficient_samples
  let variance = function
    | { n = 0 }
    | { n = 1 } -> raise Insufficient_samples
    | { sigma_x; sigma_xx; n } ->
      let n = float_of_int n in
      1. /. (n -. 1.) *. (sigma_xx -. sigma_x *. sigma_x /. n)
  let sd t = sqrt (variance t)
  let u = function
    | { n = 0 } -> raise Insufficient_samples
    | { sigma_x; n } -> sigma_x /. (float_of_int n)
  let add { sigma_x; sigma_xx; n } x =
    let sigma_x = sigma_x +. x in
    let sigma_xx = sigma_xx +. x *. x in
    let n = n + 1 in
    { sigma_x; sigma_xx; n }
end

(* Assume a list of floats is normally-distributed and return the distribution *)
let assume_normal = List.fold_left Normal.add Normal.empty

(* Take only non-None elements i.e. ignore items we couldn't parse *)

let take_some = List.fold_left (fun acc x -> match x with None -> acc | Some x -> x :: acc) []

(* Take the nth item in a list i.e. a column *)
let nth n x = try Some (List.nth x n) with _ -> None

(* Parse a line into a list of floats, ignoring '#' comments *)
let parse_line txt =
  let data = match Stringext.cut ~on:"#" txt with
    | None -> txt
    | Some (before, _after) -> before in
  Stringext.split ~on:' ' data
  |> List.filter (fun x -> String.trim x <> "")
  |> List.map (fun x ->
    try
      float_of_string x
    with
    | e -> failwith ("Failed to parse string as float: [" ^ x ^ "] in [" ^ txt ^ "]")
  )

(* Read a file into a list of strings *)
let read_file filename =
  let results = ref [] in
  let ic = open_in filename in
  try
    while true do
      let line = input_line ic in
      results := line :: !results
    done;
    assert false (* never gets here *)
  with
  | End_of_file ->
    close_in ic;
    List.rev !results

module Run = struct
  type t = {
    filename: string; (* of the test directory *)
    uuid: string;
    timestamp: string;
    version: Version.t;
    test: string;
    technology: string;
  }
  let of_string path filename = (* the results.dat file *)
    match Stringext.split ~on:'/' filename with
    | uuid :: test :: technology :: timestamp :: "results.dat" :: [] ->
      let version = Version.of_string @@ List.hd @@ read_file (path / uuid / test / technology / timestamp / "version") in
      let filename = path / uuid / test / technology / timestamp in
      { filename; uuid; timestamp; version; test; technology }
    | _ ->
      failwith (Printf.sprintf "Configuration.of_string: %s" filename)
  let to_label { version; test; technology } =
    Printf.sprintf "%s-%s:%s" (Version.to_string version) test technology
end

module String = struct
  include String
  module Map = Map.Make(String)
  module Set = Set.Make(String)
end

(* test name -> software version -> list of runs *)
let all : Run.t list Version.Map.t String.Map.t ref = ref String.Map.empty

let find path subdir filename =
  let rec loop (visited, acc) subdir =
    List.fold_left (fun (visited, acc) x ->
      let path' = subdir / x in
      if Sys.is_directory (path / path') then begin
        if String.Set.mem path' visited
        then visited, acc
        else loop (String.Set.add path' visited, acc) path'
      end else if x = filename then (visited, path' :: acc) else visited, acc
    ) (visited, acc) (Array.to_list (Sys.readdir (path / subdir))) in
  snd @@ loop (String.Set.empty, []) subdir

let load_runs path subdir =
  (* find all result.dat files inside uuid / test *)
  let result_dats = find path subdir "results.dat" in
  List.iter
    (fun result_dat ->
      let run = Run.of_string path result_dat in
      let test =
        if String.Map.mem run.Run.test !all
        then String.Map.find run.Run.test !all
        else Version.Map.empty in
      let runs =
        if Version.Map.mem run.Run.version test
        then Version.Map.find run.Run.version test
        else [] in
      all := String.Map.add run.Run.test (Version.Map.add run.Run.version (run :: runs) test) !all
    ) result_dats

let _ =
  let macperf = Hardware.id () in
  let results_dir = Sys.argv.(1) in
  load_runs results_dir (macperf / "docker-proxy-socket");

  (* test.dat will contain one line per version, one column per blocksize *)
  String.Map.iter
    (fun test versions ->
      let all_versions = List.sort Version.compare @@ List.map fst @@ Version.Map.bindings versions in
      (* Omitted the half-baked native/unknown and the dev builds*)
      let all_versions = List.filter Version.is_interesting_build all_versions in
(*
      (* Omit the native results for now *)
      let all_versions = List.filter (fun x -> x <> Version.Native) all_versions in
*)
      (try Unix.mkdir "output" 0o0755 with Unix.Unix_error(Unix.EEXIST, _, _) -> ());
      (try Unix.mkdir "output" 0o0755 with Unix.Unix_error(Unix.EEXIST, _, _) -> ());

      let gp_oc = open_out ("output/" ^ test ^ ".gp") in
      let gp fmt = Printf.ksprintf (fun s -> output_string gp_oc s; output_string gp_oc "\n") fmt in

      let xtics = fst @@ List.fold_left (fun (acc, idx) x -> Printf.sprintf "\"%s\" %d" (Version.to_pretty_string x) idx :: acc, idx+1) ([],0) all_versions in
      gp "set terminal png";
      gp "set output 'output/%s.png'" test;
      gp "set title '%s'" test;
      gp "set ylabel 'Throughput in MiB/s'";
      gp "set xlabel 'Software version'";
      gp "set xrange [-1: %d]" (List.length all_versions + 3);
      gp "set xtics rotate by -45 (%s)" (String.concat ", " xtics);
      gp "plot 'output/%s-upload.dat' using 0:1:2:3 title 'upload' linewidth 2 with yerrorbars, \\" test;
      gp "     'output/%s-download.dat' using 0:1:2:3 title 'download' linewidth 2 with yerrorbars" test;
      close_out gp_oc;

      let up = open_out ("output/" ^ test ^ "-upload.dat") in
      let down = open_out ("output/" ^ test ^ "-download.dat") in

      Printf.fprintf up "# version mean mean-sd mean+sd\n";
      Printf.fprintf down "# version mean mean-sd mean+sd\n";
      List.iter
        (fun version ->
          Printf.fprintf up "# %s\n" (Version.to_string version);
          Printf.fprintf down "# %s\n" (Version.to_string version);
          let runs = Version.Map.find version versions in
          let take_column n =
            List.map
              (fun run ->
                Printf.fprintf stderr "processing %s\n" run.Run.filename;
                run.Run.filename / "results.dat"
                |> read_file
                |> List.map parse_line
                |> List.map (nth n) (* copy-up *)
                |> take_some (* ignore lines we can't parse e.g. comments *)
                |> List.hd
              ) runs in
          (* Time for 100 MiB, convert to MiB/sec *)
          let to_mib_per_sec = List.map (fun x -> 100. /. x) in
          let write_normal oc samples =
            try
              let d = assume_normal samples in
              let u = Normal.u d in
              let sd = Normal.sd d in
              Printf.fprintf oc "%.1f %.1f %.1f\n" u (u -. sd) (u +. sd);
            with _ ->
              Printf.fprintf stderr "Skipping version %s: only %d samples\n" (Version.to_string version) (List.length samples) in

          take_column 2 |> to_mib_per_sec |> write_normal up;
          take_column 3 |> to_mib_per_sec |> write_normal down;
        ) all_versions;
      close_out up;
      close_out down;
    ) !all
