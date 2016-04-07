<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/adminguide/"]
title = "Troubleshoot DTR"
description = "Learn how to troubleshoot your DTR installation."
keywords = ["docker, registry, monitor, troubleshoot"]
[menu.main]
parent="dtr_menu_monitor_troubleshoot"
weight=10
+++
<![end-metadata]-->

# Troubleshoot DTR


## Emergency access to the Trusted Registry

If your authenticated or public access to the Trusted Registry UI has stopped
working, but your Trusted Registry admin container is still running, you can add
an
[ambassador container](https://docs.docker.com/articles/ambassador_pattern_linking/)
to get temporary unsecure access to it.

For Trusted Registry version 1.4.3, run the following command in a Trusted Registry CLI:

```
docker run --rm -it --net dtr -p 9999:80 svendowideit/ambassador dockertrustedregistry_admin_server_1 80
```
However, if you are running a version prior to it,  1.4.2 or earlier, then continue to run this command:

```
$ docker run --rm -it --link docker_trusted_registry_admin_server:admin -p 9999:80 svendowideit/ambassador
```

Either command gives you access on port `9999` on your Trusted Registry server
`http://<dtr-host-ip>:9999`. This guide assumes that you are a member of the `docker` group, or you  have root privileges. Otherwise, you may need to add `sudo` to the previous example command.

### SSH access to host

As an extra measure of safety, ensure you have SSH access to the Trusted
Registry host before you start using it.

If you are hosting Trusted Registry on an EC2 host launched from the AWS
Marketplace AMI, note that the user is `ec2-user`:
`/path/to/private_key/id_rsa ec2-user@<dtr-dns-entry>`.


## Client Docker Daemon diagnostics

To debug client Docker daemon communication issues with the Trusted Registry,
Docker also provides a diagnostics tool to be run on the client Docker daemon.

> **Warning:** These diagnostics files may contain secrets that you need to remove before passing on, such as raw container log files, Azure storage credentials, or passwords that may be sent to non-Docker Trusted Registry containers using the `docker run -e PASSWORD=asdf` environment variable options.

If you supply an administrator username and password, then the `diagnostics`
tool also downloads additional logs and configuration data from the remote
Trusted Registry server. Download and run this tool using the following command:

```
$ wget https://dhe.mycompany.com/admin/bin/diagnostics && chmod +x diagnostics
$ sudo ./diagnostics dhe.mycompany.com > enduserDiagnostics.zip DTR
administrator password (provide empty string if there is no admin server
authentication):
WARN  [1.1.0-alpha-001472_g8a9ddb4] Encountered errors running diagnostics
errors=[Failed to copy DTR Adminserver's exported settings into ZIP output:
"Failed to read next tar header: \"archive/tar: invalid tar header\"" Failed to
copy logs from DTR Adminserver into ZIP output: "Failed to read next tar header:
\"archive/tar: invalid tar header\"" error running "sestatus": "exit status 127"
error running "dmidecode": "exit status 127"]
```

The zip file contains the following information:

- your local Docker host's `ca-certificates.crt`
- `containers/`: the first 20 running, stopped and paused containers `docker inspect`
  information and log files.
- `dockerEngine/`: the local Docker daemon's `info` and `version` output
- `dockerState/`: the local Docker daemon's container states, image states, log file, and daemon configuration file
- `dtr/`: Remote Docker Trusted Registry services information. This directory will only be populated if the user enters a Docker Trusted Registry "admin" username and password.
- - `dtr/logs/`: the remote Docker Trusted Registry container log files. This directory will only be populated if the user enters a Docker Trusted Registry "admin" username and password.
- - `dtr/exportedSettings/`: the Docker Trusted Registry manager container's log files and a backup of the `/usr/local/etc/dtr` Docker Trusted Registry configuration directory. See the [export settings section](#export-settings) for more details.
- `sysinfo/`: local Host information
- `errors.txt`: errors and warnings encountered while running diagnostics

### Starting and stopping the Trusted Registry

If you need to stop and/or start the Trusted Registry (for example, upgrading, or troubleshooting), use the following commands:

`sudo bash -c "$(docker run docker/trusted-registry stop)"`


`sudo bash -c "$(docker run docker/trusted-registry start)"`
