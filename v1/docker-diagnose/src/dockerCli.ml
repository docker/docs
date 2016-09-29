open Common
open Misc

let name = "docker-cli"
let docker_sock = home / "Library/Containers/com.docker.docker/Data/s60"

let check_version expected_version =
  match Cmd.docker ?path:app_location "-v" with
  | []   -> error name "empty version"
  | v::_ ->
    match Stringext.cut ~on:"," v with
    | None -> error name "%s is not a valid version" v
    | Some (v, _) ->
      let v = List.hd (List.rev (Stringext.split ~on:' ' v)) in
      if v <> expected_version then
        error name "wrong Docker version: %s (expected %s)" v expected_version

let check () =
  let docker = match app_location with
    | None   -> "docker"
    | Some p -> p / "Contents/Resources/bin/docker"
  in
  check_cmd name docker;
  check_version "1.12.2-rc1";
  check_socket name "/var/run/docker.sock";

  (* s60: formerly /var/tmp/docker.sock *)
  check_socket name docker_sock;
  begin
    try
      ignore(Cmd.docker "ps")
    with e ->
      error ~e "docker-cli" "docker ps failed"
  end;
  ok name
