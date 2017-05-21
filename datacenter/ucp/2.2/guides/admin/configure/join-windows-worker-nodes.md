---
title: Join Windows worker nodes to a swarm
description: Join worker nodes that are running on Windows Server 2016 to a swarm managed by UCP. 
keywords: UCP, swarm, Windows, cluster
---

UCP supports worker nodes that run on Windows Server 2016. Only worker nodes
are supported on Windows, and all manager nodes in the swarm must run on Linux.

Follow these steps to enable a worker node on Windows.

1.  Install UCP on a Linux distribution.
2.  Install Docker EE on Windows Server 2016.
3.  Configure the Windows node.
4.  Join the Windows node to the swarm.  

## Install UCP

To install UCP, follow the instructions in 
[Learn how to install Docker Universal Control Plane on production](/datacenter/ucp/2.2/guides/admin/install/).
UCP requires Docker EE version 17.06 or later.

For internal development, you need to `docker login` and pull
the UCP images manually. For Beta, the images will be available publicly.

## Install Docker EE on Windows Server 2016

Install Docker Enterprise Edition (*Docker EE*) on a Windows
Server 2016 instance to enable joining a swarm that's managed by UCP.
Follow the instructions in 
[Install Docker Enterprise Edition for Windows Server 2016](/docker-ee-for-windows/install/#using-a-script-to-install-docker-ee).
For internal development, install the dev binaries in the zip archive at
[windows/amd64/docker-17.06.0-dev.zip](https://master.dockerproject.org/windows/amd64/docker-17.06.0-dev.zip),
because you need version 17.06 or later to join a UCP swarm. For Beta, the binaries 
will be available publicly at [download.docker.com](https://download.docker.com/components/engine/windows-server).

## Configure the Windows node

Follow these steps to configure the docker daemon and the Windows environment to enable
joining a UCP cluster.

1.  Pull the Windows-specific image of `ucp-agent`, which is named `ucp-agent-win`.
2.  Run the Windows worker setup script provided with `ucp-agent-win`.
3.  Join the swarm with the token provided by the UCP web UI. 

### Pull the Windows-specific images

On a manager node, run the following command to list the images that are required
on Windows nodes.

```bash
$ docker run --rm dockerorcadev/ucp:2.2.0-latest images --list --image-version dev: --enable-windows
dockerorcadev/ucp-agent-win:2.2.0-5213679
dockerorcadev/ucp-dsinfo-win:2.2.0-5213679
```

On Windows Server 2016, in a PowerShell terminal running as Administrator,
login to Docker Hub with the `docker login` command and pull the listed images. 

```ps
PS> docker pull dockerorcadev/ucp-agent-win:2.2.0-5213679
```

### Run the Windows node setup script

You need to open ports 2376 and 12376, and create certificates
for the Docker daemon to communicate securely. Run this command:

```ps
PS> docker run --rm dockerorcadev/ucp-agent-win:2.2.0-5213679 windows-script | powershell -noprofile -noninteractive -command 'Invoke-Expression -Command $input'
```

The Windows node is ready to join the swarm.

## Join the Windows node to the swarm

Now you can join the UCP cluster by using the `docker swarm join` command that's
provided by the UCP web UI. For more info, see
[Join worker nodes](/datacenter/ucp/2.2/guides/admin/install/#step-7-join-worker-nodes).
The command looks similar to the following.

```ps
PS> docker swarm join --token <token> <ucp-manager-ip>
```

Run the setup script on each instance of Windows Server that will be a worker node.

## Configure a Windows worker node manually  

The following sections describe how to run the commands in the setup script
manually to configure the `dockerd` service and the Windows environment.
The script opens ports in the firewall and sets up certificates for `dockerd`.

>**Note**: For internal development, you need to run these commands manually, 
because the script assumes access to public images. You need to be logged in to Docker Hub.

To see the script, you can run the `windows-script` command without piping
to the `Invoke-Expression` cmdlet.

```ps
PS> docker run --rm dockerorcadev/ucp-agent-win:2.2.0-latest windows-script
```

### Open ports in the Windows firewall

UCP and Docker EE require that ports 2376 and 12376 are open for inbound
TCP traffic.

In a PowerShell terminal running as Administrator, run these commands
to add rules to the Windows firewall.

```ps
PS> netsh advfirewall firewall add rule name="docker_local" dir=in action=allow protocol=TCP localport=2376
PS> netsh advfirewall firewall add rule name="docker_proxy" dir=in action=allow protocol=TCP localport=12376
```

###  Set up certs for the dockerd service

1.  Create the directory `C:\ProgramData\docker\daemoncerts`.
2.  In a PowerShell terminal running as Administrator, run the following command
    to generate certificates. 
    ```ps
    PS> docker run --rm -v C:\ProgramData\docker\daemoncerts:C:\certs dockerorcadev/ucp-agent-win:2.2.0-5213679 generate-certs
    ```
3.  To set up certificates, run the following commands to stop and unregister the
    `dockerd` service, register the service with the certificates, and restart the service.

    ```ps
    PS> Stop-Service docker
    PS> dockerd --unregister-service
    PS> dockerd -H npipe:// -H 0.0.0.0:2376 --tlsverify --tlscacert=C:\ProgramData\docker\daemoncerts\ca.pem --tlscert=C:\ProgramData\docker\daemoncerts\cert.pem --tlskey=C:\ProgramData\docker\daemoncerts\key.pem --register-service
    PS> Start-Service docker
    ```

The `dockerd` service and the Windows environment are now configured to join a UCP swarm.

>**Tip:** If the TLS certificates aren't set up correctly, the UCP web UI shows the
>following warning.

```
Node WIN-NOOQV2PJGTE is a Windows node that cannot connect to its local Docker daemon.
```