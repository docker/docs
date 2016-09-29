open Misc
open Common

let check_vtx () =
  let name = "VT-x" in
  let features = List.tl @@ Stringext.split ?max:None ~on:' ' @@ String.concat " " @@ Cmd.read_stdout "sysctl -a | grep machdep.cpu.features" in
  if List.mem "VMX" features then ok name else error name "No VMX in machdep.cpu.features: %s" (String.concat " " features)

let check_kern_hv_support () =
  let name = "kern.hv_support" in
  try
    Cmd.exec "sysctl -a | grep 'kern.hv_support: 1'";
    ok name
  with _ ->
    error name "sysctl -a reports no hypervisor support"
