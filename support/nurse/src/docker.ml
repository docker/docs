open Lwt.Infix

let docker = "/usr/local/bin/docker"

let inspect name = Command.run docker [ "inspect"; name ]

let run args = Command.run docker ("run" :: args)
