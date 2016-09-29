open Common

let name = "menubar"
let check () =
  check_no_lctl name "com.docker.docker-menu";
  check_no_ps name "com.docker.docker-menu";
  ok name
