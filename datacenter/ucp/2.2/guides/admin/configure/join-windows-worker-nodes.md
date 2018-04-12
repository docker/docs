---
title: Join Windows worker nodes to a swarm
description: Join worker nodes that are running on Windows Server 2016 to a swarm managed by UCP.
keywords: UCP, swarm, Windows, cluster
---

UCP supports worker nodes that run on Windows Server 2016. Only worker nodes
are supported on Windows, and all manager nodes in the swarm must run on Linux.

Follow these steps to enable a worker node on Windows.

1.  Install UCP on a Linux distribution.
2.  Install Docker Enterprise Edition (*Docker EE*) on Windows Server 2016.
3.  Configure the Windows node.
4.  Join the Windows node to the swarm.

## Install UCP

Install UCP on a Linux distribution.
[Learn how to install UCP on production](../install/index.md).
UCP requires Docker EE version 17.06 or later.

## Install Docker EE on Windows Server 2016

[Install Docker EE](/docker-ee-for-windows/install/#using-a-script-to-install-docker-ee)
on a Windows Server 2016 instance to enable joining a swarm that's managed by
UCP.

## Configure the Windows node

Follow these steps to configure the docker daemon and the Windows environment.

1.  Pull the Windows-specific image of `ucp-agent`, which is named `ucp-agent-win`.
2.  Run the Windows worker setup script provided with `ucp-agent-win`.
3.  Join the swarm with the token provided by the UCP web UI.

### Pull the Windows-specific images

On a manager node, run the following command to list the images that are required
on Windows nodes.

```bash
docker container run --rm {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} images --list --enable-windows \
{{ page.ucp_org }}/ucp-agent-win:{{ page.ucp_version }} \
{{ page.ucp_org }}/ucp-dsinfo-win:{{ page.ucp_version }}
```

On Windows Server 2016, in a PowerShell terminal running as Administrator,
log in to Docker Hub with the `docker login` command and pull the listed images.

```powershell
docker image pull {{ page.ucp_org }}/ucp-agent-win:{{ page.ucp_version }}
docker image pull {{ page.ucp_org }}/ucp-dsinfo-win:{{ page.ucp_version }}
```

### Run the Windows node setup script

You need to open ports 2376 and 12376, and create certificates
for the Docker daemon to communicate securely. Run this command:

```powershell
$script = [ScriptBlock]::Create((docker run --rm {{ page.ucp_org }}/ucp-agent-win:{{ page.ucp_version }} windows-script | Out-String))

Invoke-Command $script
```

> Docker daemon restart
>
> When you run `windows-script`, the Docker service is unavailable temporarily.

The Windows node is ready to join the swarm. Run the setup script on each
instance of Windows Server that will be a worker node.

### Compatibility with daemon.json

The script may be incompatible with installations that use a config file at
`C:\ProgramData\docker\config\daemon.json`. If you use such a file, make sure
that the daemon runs on port 2376 and that it uses certificates located in
`C:\ProgramData\docker\daemoncerts`. If certificates don't exist in this
directory, run `ucp-agent-win generate-certs`, as shown in Step 2 of the
[Set up certs for the dockerd service](#set-up-certs-for-the-dockerd-service)
procedure.

In the daemon.json file, set the `tlscacert`, `tlscert`, and `tlskey` options
to the corresponding files in `C:\ProgramData\docker\daemoncerts`:

```json
{
...
		"debug":     true,
		"tls":       true,
		"tlscacert": "C:\ProgramData\docker\daemoncerts\ca.pem",
		"tlscert":   "C:\ProgramData\docker\daemoncerts\cert.pem",
		"tlskey":    "C:\ProgramData\docker\daemoncerts\key.pem",
		"tlsverify": true,
...
}
```

## Join the Windows node to the swarm

Now you can join the UCP cluster by using the `docker swarm join` command that's
provided by the UCP web UI. [Learn to add nodes to your swarm](scale-your-cluster.md).
The command looks similar to the following.

```powershell
docker swarm join --token <token> <ucp-manager-ip>
```

Run the `docker swarm join` command on each instance of Windows Server that
will be a worker node.

## Configure a Windows worker node manually

The following sections describe how to run the commands in the setup script
manually to configure the `dockerd` service and the Windows environment.
The script opens ports in the firewall and sets up certificates for `dockerd`.

To see the script, you can run the `windows-script` command without piping
to the `Invoke-Expression` cmdlet.

```powershell
docker container run --rm {{ page.ucp_org }}/ucp-agent-win:{{ page.ucp_version }} windows-script
```


### Open ports in the Windows firewall

UCP and Docker EE require that ports 2376 and 12376 are open for inbound
TCP traffic.

In a PowerShell terminal running as Administrator, run these commands
to add rules to the Windows firewall.

```powershell
netsh advfirewall firewall add rule name="docker_local" dir=in action=allow protocol=TCP localport=2376
netsh advfirewall firewall add rule name="docker_proxy" dir=in action=allow protocol=TCP localport=12376
```

### Set up certs for the dockerd service

1.  Create the directory `C:\ProgramData\docker\daemoncerts`.
2.  In a PowerShell terminal running as Administrator, run the following command
    to generate certificates.

    ```powershell
    docker container run --rm -v C:\ProgramData\docker\daemoncerts:C:\certs {{ page.ucp_org }}/ucp-agent-win:{{ page.ucp_version }} generate-certs
    ```

3.  To set up certificates, run the following commands to stop and unregister the
    `dockerd` service, register the service with the certificates, and restart the service.

    ```powershell
    Stop-Service docker
    dockerd --unregister-service
    dockerd -H npipe:// -H 0.0.0.0:2376 --tlsverify --tlscacert=C:\ProgramData\docker\daemoncerts\ca.pem --tlscert=C:\ProgramData\docker\daemoncerts\cert.pem --tlskey=C:\ProgramData\docker\daemoncerts\key.pem --register-service
    Start-Service docker
    ```

The `dockerd` service and the Windows environment are now configured to join a UCP swarm.

> **Tip:** If the TLS certificates aren't set up correctly, the UCP web UI shows the
> following warning.

```
Node WIN-NOOQV2PJGTE is a Windows node that cannot connect to its local Docker daemon.
```

## Uninstall UCP from Windows Server

The following steps return the Docker Engine to its original configuration:

1. Unregister the docker service and register it again without the TLS
   certificates:

   ```powershell
   Stop-Service docker
   dockerd --unregister-service
   dockerd -H npipe:// --register-service
   Start-Service docker
   ```

2. Remove the `certs` directory for the docker service:

   ```powershell
   Remove-Item -Recurse C:\ProgramData\docker\daemoncerts
   ```

3. Remove the firewall rules:

   ```powershell
   netsh advfirewall firewall delete rule name="docker_2376_in"
   netsh advfirewall firewall delete rule name="docker_12376_in"
   netsh advfirewall firewall delete rule name="docker_2377_in"
   netsh advfirewall firewall delete rule name="docker_4789_in"
   netsh advfirewall firewall delete rule name="docker_4789_out"
   netsh advfirewall firewall delete rule name="docker_7946_in"
   netsh advfirewall firewall delete rule name="docker_7946_out"
   ```

## Windows nodes limitations

Some features are not yet supported on Windows nodes:

* Networking
  * The swarm mode routing mesh can't be used on Windows nodes. You can expose
  a port for your service in the host where it is running, and use the HTTP
  routing mesh to make your service accessible using a domain name.
  * Encrypted networks are not supported. If you've upgraded from a previous
  version, you also need to recreate the `ucp-hrm` network to make it
  unencrypted.
* Secrets
  * When using secrets with Windows services, Windows stores temporary secret
  files on disk. You can use BitLocker on the volume containing the Docker
  root directory to encrypt the secret data at rest.
  * When creating a service which uses Windows containers, the options to
  specify UID, GID, and mode are not supported for secrets. Secrets are
  currently only accessible by administrators and users with system access
  within the container.
* Mounts
  * On Windows, Docker can't listen on a Unix socket. Use TCP or a named pipe
  instead.
