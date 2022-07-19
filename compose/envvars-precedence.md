---
title: Environment variables precedence
description: TODO
keywords: compose, environment, env file
---

<style>
div {
  color: RoyalBlue; 
}
table {
  font-family: monospace;
}
</style>

- Command Line (docker compose run --env <KEY[=VAL]> https://docs.docker.com/engine/reference/commandline/compose_run/#options)
- Compose File (service::environment section: https://docs.docker.com/compose/compose-file/#environment)
- Compose File (service::env_file section file: https://docs.docker.com/compose/compose-file/#env_file)
- Container Image ENV directive (https://docs.docker.com/engine/reference/builder/#env)

When `WHEREAMI` is the variable in case:

| Image (`ENV` directive in Dockerfile) |  OS Environment   | `.env` file on the project root (or overwrite through `docker compose -–env-file <FILE>`) | Compose file (`service::env_file`) |   Compose file (`service::environment`)   | Command line (docker compose run -e <KEY[=VAL]>) |          RESULT          |
|:-------------------------------------:|:-----------------:|:-----------------------------------------------------------------------------------------:|:----------------------------------:|:-----------------------------------------:|:------------------------------------------------:|:------------------------:|
|         <div>Dockerfile</div>         |      OS_Env       |                                   WHEREAMI=DotEnv_File                                    |                                    |                                           |                                                  |        Dockerfile        |
|              Dockerfile               |      OS_Env       |                              WHEREAMI=<div>DotEnv_File</div>                              |                                    |                                           |                     WHEREAMI                     |       DotEnv_File        |
|              Dockerfile               | <div>OS_Env</div> |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |                                           |                     WHEREAMI                     |          OS_Env          |
|         <div>Dockerfile</div>         |      OS_Env       |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |                                           |                                                  |        Dockerfile        |
|              Dockerfile               |                   |                 WHEREAMI=${WHEREAMI:-<div>DotEnv_File_DefaultValue</div>}                 |                                    |                                           |                     WHEREAMI                     | DotEnv_File_DefaultValue |
|              Dockerfile               |      OS_Env       |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |                                           |         WHEREAMI=<div>Command_Line</div>         |       Command_Line       |
|              Dockerfile               | <div>OS_Env</div> |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |                 WHEREAMI                  |                                                  |          OS_Env          |
|              Dockerfile               |      OS_Env       |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |  WHEREAMI=<div>Environment_Section</div>  |                                                  |   Environment_Section    |
|              Dockerfile               |      OS_Env       |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |       WHEREAMI=Environment_Section        |         WHEREAMI=<div>Command_Line</div>         |       Command_Line       |
|              Dockerfile               | <div>OS_Env</div> |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |       WHEREAMI=Environment_Section        |                     WHEREAMI                     |          OS_Env          |
|              Dockerfile               |      OS_Env       |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    |                 WHEREAMI                  |         WHEREAMI=<div>Command_Line</div>         |       Command_Line       |
|              Dockerfile               |      OS_Env       |                                                                                           |         WHEREAMI=Env_File          |                                           |         WHEREAMI=<div>Command_Line</div>         |       Command_Line       |
|              Dockerfile               |      OS_Env       |                                                                                           |         WHEREAMI=Env_File          |  WHEREAMI=<div>Environment_Section</div>  |                                                  |   Environment_Section    |
|              Dockerfile               |      OS_Env       |                                                                                           |    WHEREAMI=<div>Env_File</div>    |                                           |                                                  |         Env_File         |
|              Dockerfile               | <div>OS_Env</div> |                                                                                           |                                    | WHEREAMI=${WHEREAMI:-Environment_Section} |                                                  |          OS_Env          |
|              Dockerfile               | <div>OS_Env</div> |                      WHEREAMI=${WHEREAMI:-DotEnv_File_DefaultValue}                       |                                    | WHEREAMI=${WHEREAMI:-Environment_Section} |                                                  |          OS_Env          |
