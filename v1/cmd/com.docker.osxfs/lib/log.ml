
let name = "osxfs"
let doc = "host filesharing daemon"

let src =
  let src = Logs.Src.create name ~doc in
  Logs.Src.set_level src (Some Logs.Debug);
  src

module Log = (val Logs.src_log src : Logs.LOG)

let src_info =
  let src = Logs.Src.create name ~doc in
  Logs.Src.set_level src (Some Logs.Info);
  src

module LogInfo = (val Logs.src_log src_info : Logs.LOG)
