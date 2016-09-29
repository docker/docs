open Misc

(* Categories of data from system_profiler *)
let interesting = [
  "SPFirewallDataType"; (* firewall or not *)
  "SPHardwareDataType"; (* hw serial number, cores, boot ROM *)
]

let collect () =
  List.iter (fun category ->
    let name = Filename.concat Logs.dir category in
    Cmd.exec "system_profiler %s > %s" category name;
  ) interesting
