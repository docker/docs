#! /usr/bin/env ocamlscript
Ocaml.packs := [ "unix"; "stringext" ]
--

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
  let numbers = Stringext.split ~on:' ' data in
  List.map float_of_string numbers

let setify = List.fold_left (fun set x -> if List.mem x set then set else x :: set) []

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
  | End_of_file -> List.rev !results

module Configuration = struct
  type t = {
    uuid: string;
    direction: string;
    technology: string;
  }
  let of_string filename =
    match Stringext.split ~on:'.' (Filename.basename filename) with
    | uuid :: direction :: technology :: "dat" :: [] ->
      { uuid; direction; technology }
    | _ ->
      failwith (Printf.sprintf "Configuration.of_string: %s" filename)
  let to_label { uuid; direction; technology } =
    Printf.sprintf "%s:%s:%s" (String.sub uuid 0 4) direction technology
end

let _ =
  let files = List.tl @@ Array.to_list Sys.argv in
  let distributions = List.map (fun filename ->
    filename
    |> read_file
    |> List.map parse_line
    |> List.map (nth 1) (* ignore timestamp *)
    |> take_some (* ignore lines we can't parse e.g. comments *)
    |> assume_normal
  ) files in
  let all = List.combine (List.map Configuration.of_string files) distributions in
  (* Leave out those with no samples *)
  let all = List.filter (fun (_, d) -> d.Normal.n > 0) all in

  (* Only include a hardware id if we have all 4 columns *)
  let all_uuids = setify @@ List.map (fun (config, _) -> config.Configuration.uuid) all in
  let include_uuids = List.filter (fun uuid -> List.length (List.filter (fun (config, _) -> config.Configuration.uuid = uuid) all) = 4) all_uuids in
  let all = List.filter (fun (config, _) -> List.mem config.Configuration.uuid include_uuids) all in

  (* Give them all an index *)
  let _, all = List.fold_left (fun (i, acc) x -> i + 1, (i, x) :: acc) (1, []) all in

  let gp = open_out "output/plot.gp" in
  Printf.fprintf gp "
  set terminal png
  set output 'output/plot.png'
  set title 'iperf throughput comparison'
  set ylabel 'Throughput in Mbit/s'
  set xlabel 'Hardware uuid : to/from container : hyperkit+vmnet/vbox'
  set nokey
  ";
  let xtics = List.map (fun (i, (config, _)) -> "\""^(Configuration.to_label config)^"\" " ^ (string_of_int i)) all in
  Printf.fprintf gp "
  set xrange [0: %d]
  set xtics rotate by -45 (%s)
  plot 'output/plot.normal' linestyle 2 linewidth 2 with yerrorbars
  " (List.length all + 1) (String.concat ", " xtics);
  close_out gp;

  let dat = open_out "output/plot.normal" in
  Printf.fprintf dat "# 'bar' 'mean' 'mean-sd' 'mean+sd'\n";
  List.iter (fun (i, (config, distribution)) ->
    Printf.fprintf dat "# %s\n%!" (Configuration.to_label config);
    try
      let u = Normal.u distribution in
      let sd = Normal.sd distribution in
      Printf.fprintf dat "%d %.1f %.1f %.1f\n" i u (u -. sd) (u +. sd);
    with Normal.Insufficient_samples ->
      Printf.fprintf dat "# insufficient samples\n"
  ) all;
  close_out dat
