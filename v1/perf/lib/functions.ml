let x = 1

module Version = struct
  module M = struct
    type t =
      | Pinata of int list
      | PinataDev of int option * string
      | Virtualbox
      | Fusion
      | Native

    let of_string x =
      try
          if String.lowercase x = "virtualbox"
          then Virtualbox
          else if String.lowercase x = "vmwarefusion"
          then Fusion
          else if String.lowercase x = "unknown"
          then Native
          else match Stringext.split ~on:'.' x with
          | "dev" :: nr :: desc -> begin match (try Some (int_of_string nr) with _ -> None) with
	    | Some sort_order -> PinataDev (Some sort_order, (String.concat "." desc))
	    | None -> PinataDev (None, (String.concat "." (nr::desc))); end
          | vers -> Pinata (List.map int_of_string vers)
      with
      | _ -> failwith (Printf.sprintf "Failed to parse version [%s]" x)
    let to_string = function
      | Virtualbox -> "virtualbox"
      | Fusion -> "vmwarefusion"
      | Native -> "unknown"
      | Pinata x -> String.concat "." @@ List.map string_of_int x
      | PinataDev (_,x) -> x

    let alpha7 = [1;10;0;2338]
    let alpha8 = [1;10;0;2536]
    let alpha9 = [1;10;0;2703]
    let alpha10 = [1;10;0;2898]
    let alpha11 = [1;10;0;3070]
    let alpha12 = [1;10;2;3360]
    let beta1 = [1;10;2;3791]
    let beta2 = [4054]
    let beta3 = [4433]
    let beta4 = [4782]
    let beta5 = [5049]
    let beta6 = [5404]
    let beta7 = [5830]
    let beta8 = [6072]
    let beta9 = [6388]
    let beta10 = [6662]
    let beta11 = [6974]
    let released_builds = [ alpha7; alpha8; alpha9; alpha10; alpha11; alpha12;
      beta1; beta2; beta3; beta4; beta5; beta6; beta7; beta8; beta9; beta10; beta11 ]

    let is_released_build = function
      | Virtualbox
      | Fusion -> true
      | Pinata x -> List.mem x released_builds
      | _ -> false

    let is_interesting_build = function
      | PinataDev (_,_) -> true
      | x -> is_released_build x

    let to_pretty_string = function
      | Virtualbox -> "Virtualbox 5"
      | Fusion -> "VMware Fusion 8"
      | Native -> "Bare metal"
      | Pinata [1;10;0;2338] -> "Alpha 7"
      | Pinata [1;10;0;2536] -> "Alpha 8"
      | Pinata [1;10;0;2703] -> "Alpha 9"
      | Pinata [1;10;0;2898] -> "Alpha 10"
      | Pinata [1;10;0;3070] -> "Alpha 11"
      | Pinata [1;10;2;3360] -> "Alpha 12"
      | Pinata [1;10;2;3791] -> "Beta 1"
      | Pinata [4054] -> "Beta 2"
      | Pinata [4433] -> "Beta 3"
      | Pinata [4782] -> "Beta 4"
      | Pinata [5049] -> "Beta 5"
      | Pinata [5404] -> "Beta 6"
      | Pinata [5830] -> "Beta 7"
      | Pinata [6072] -> "Beta 8"
      | Pinata [6388] -> "Beta 9"
      | Pinata [6662] -> "Beta 10"
      | Pinata [6974] -> "Beta 11"
      | Pinata x -> "build " ^ (String.concat "." @@ List.map string_of_int x)
      | PinataDev (_,x) -> "dev build " ^ x

    let compare x y =
      (* compare based on constructor, and then version number *)
      let rec ints x y = match x, y with
        | [], [] -> 0
        | [], _::_ -> 1
        | _::_, [] -> -1
        | x :: xs, y :: ys ->
          let c = compare x y in
          if c = 0 then ints xs ys else c in
      let to_int = function
        | Pinata _ -> 0
        | PinataDev _ -> 1
        | Virtualbox -> 2
        | Fusion -> 3
        | Native -> 4 in
      match x, y with
        | Pinata x, Pinata y -> ints x y
        | PinataDev (Some xo, xs), PinataDev (Some yo, ys) -> begin
	  if xo = yo then compare xs ys else compare xo yo
	end
        | PinataDev (Some _, _), PinataDev (None, _) -> -1
        | PinataDev (None, _), PinataDev (Some _, _) -> 1
        | PinataDev (None, x), PinataDev (None, y) -> compare x y
        | Virtualbox, Virtualbox -> 0
        | Fusion, Fusion -> 0
        | Native, Native -> 0
        | a, b -> compare (to_int a) (to_int b)
  end
  include M
  module Map = Map.Make(M)
  module Set = Set.Make(M)
end

module Hardware = struct
  let id () =
  try Sys.getenv "PINATA_PERF_HWID"
  with _ -> "EC420ECC-EDEE-50AC-9B9C-5984E7F4E23C"
end
