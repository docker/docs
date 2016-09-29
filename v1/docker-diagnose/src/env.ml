open Common
open Misc

let name = "environment"
let env =
  let lines = Cmd.read_stdout "env" in
  let kvs =
    List.fold_left (fun acc s ->
        if s = "" then acc
        else match Stringext.cut s ~on:"=" with
          | None       -> acc
          | Some (k,v) -> (k ,v) :: acc
      ) [] lines
  in
  let vars = List.filter (fun (k, _) ->
      String.length k > 6 && String.sub k 0 6 = "DOCKER"
    ) kvs
  in
  vars

let check () =
  let keys = List.map fst env in
  match keys with
  | []  -> debug "no DOCKER_* environment in the shell: good!"; ok name
  | [k] -> error name "the variable %s should not be set" k
  | ks  ->
    let vars = String.concat " " ks in
    error name "the variables %s should not be set" vars
