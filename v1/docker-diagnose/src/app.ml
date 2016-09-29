open Common

let name = "Docker"

let check () =
  try
    check_ps name name;
    ok name
  with
  | e ->
    error ~e name "Docker app is not running"

let collect () =
  Logs.add_string "version" Version.git
