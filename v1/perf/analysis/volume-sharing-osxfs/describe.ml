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
(*
let take_some = List.fold_left (fun acc x -> match x with None -> acc | Some x -> x :: acc) []
*)
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
    | uuid :: test :: technology :: timestamp :: "logs" :: "512.time" :: [] ->
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
  let result_dats = find path subdir "512.time" in
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
  let tests = [ "volume-sharing-dd-read"; "volume-sharing-dd-write" ] in
  let macperf = Hardware.id () in
  let results_dir = Sys.argv.(1) in
  List.iter (fun x -> load_runs results_dir (macperf / x)) tests;

  (* test.dat will contain one line per version, one column per blocksize *)
  String.Map.iter
    (fun test versions ->
      let all_versions = List.sort Version.compare @@ List.map fst @@ Version.Map.bindings versions in
      (* Omitted the half-baked native/unknown and the dev builds*)
      let all_versions = List.filter Version.is_interesting_build all_versions in

      (try Unix.mkdir "output" 0o0755 with Unix.Unix_error(Unix.EEXIST, _, _) -> ());
      let dat = open_out ("output/" ^ test ^ ".dat") in
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
      gp "plot 'output/%s.dat' using 0:1 title '512' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:2 title '1K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:3 title '2K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:4 title '4K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:5 title '8K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:6 title '16K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:7 title '32K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:8 title '64K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:9 title '128K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:10 title '256K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:11 title '512K' linewidth 2 with points, \\" test;
      gp "     'output/%s.dat' using 0:12 title '1M' linewidth 2 with points" test;

      (* Omit the native results for now *)
      let all_versions = List.filter (fun x -> x <> Version.Native) all_versions in
      Printf.fprintf dat "# version blocksize1 blocksize2 ... blocksizeN\n";
      List.iter
        (fun version ->
          Printf.fprintf dat "# %s\n" (Version.to_string version);
          let runs = Version.Map.find version versions in
          let block_sizes = Hashtbl.create 7 in
          List.iter
            (fun run ->
              Printf.fprintf stderr "processing %s\n" run.Run.filename;
              let columns =
                [ List.map
                  (fun block_size ->
                    read_file (run.Run.filename / "logs" / (Printf.sprintf "%d.time" block_size))
                    |> List.hd
                    |> String.trim
                    |> float_of_string
                  ) [ 512; 1024; 2048; 4096; 8192; 16384; 32768; 65536; 131072;
                      262144; 524288; 1048576 ] ] in
              let ncols = List.fold_left min max_int (List.map List.length columns) in
              if ncols <> max_int then begin
                for i = 0 to ncols - 1 do
                  let distribution =
                    if Hashtbl.mem block_sizes i
                    then Hashtbl.find block_sizes i
                    else
                      let d = Normal.empty in
                      Hashtbl.replace block_sizes i d;
                      d in
                  let column = List.map (fun x -> List.nth x i) columns in
                  (* Convert from time to copy 128MiB to MBytes/sec *)
                  let column = List.map (fun time -> 128.0 /. time) column in
                  let distribution' = List.fold_left Normal.add distribution column in
                  Hashtbl.replace block_sizes i distribution';
                done;
              end
            ) runs;
          let idxs = Hashtbl.fold (fun idx _ acc -> idx :: acc) block_sizes [] in
          List.iter
            (fun idx ->
              let d = Hashtbl.find block_sizes idx in
              Printf.fprintf dat "%f " (Normal.u d);
            ) (List.sort compare idxs);
          Printf.fprintf dat "\n"
        ) all_versions;
      close_out dat;
      close_out gp_oc
    ) !all
