Our goal is to provide an OSX native, seamless experience using Docker on the
Mac and on Windows. On Linux Docker is fully-integrated and everything
“just works”. Today on the Mac and Windows we have to rely on a VirtualBox or
VMware fusion VM which

- exposes a hypervisor product and a Linux VM and all the settings and state to
  the user
- creates additional failure modes:
  - VM IP address changed by the hypervisor, we lose contact
  - host DNS resolver is exposed to the VM by virtualbox, laptop joins a new
    network, DNS is broken
  - we rely on NAT which can be disabled by VPN software on the host
    (e.g. Cisco AnyConnect)
  - file ownership is different on the host vs the VM, causing failures with
    “-v /path:/path”
  - no inotify events in the VM, breaking workflows (e.g. jekyll rebuilding a blog)
  - port space on the VM is different from the host, breaking “-p”

Architecture
============

We have: (some of this is experimental!) (in no particular order)
- an embedded “hyperkit” VM to run Linux. This is completely under our control and
  we can do things like: add new types of virtual hardware (e.g. for shared
  memory communication or new networking modes) and hide it from the user
  (i.e. it’s just one of our subprocesses, there’s no separate GUI)
- a custom Alpine-based Linux distribution, optimised for size, speed of boot and
  running Docker
- a custom 9P-based system for “forwarding” /var/run/docker.sock on the Mac to the
  /var/run/docker.sock in the VM. This uses the 9P protocol over shared memory
  and avoids using the network so it cannot break due to an IP address change or
  a broken TLS cert. We would like to replace this with a socket-based solution
  e.g. HyperV sockets on windows and VSOCK on Linux.
- a custom file sharing system in which the VM uses “fuse” to mount a fileserver
  running on the Mac. This allows us to manipulate the perceived uid and gid of
  a file based on the accessing pid to work-around Mac-VM uid and gid mismatches.
  This also allows us to listen for Mac FSevents and translate these into Linux
  inotify events. The server runs as the user (not root) so it naturally can
  only see the user’s files (c.f. using the Mac in-kernel NFS server which can
  see and export everything)
- the ability to forward ports from the Mac to the VM so “docker run -p 1234”
  works as expected, without having to expose the internal IP address of the VM.
- a user-space TCP/IP stack which can proxy VM TCP, UDP, DNS via the Mac’s
  native application sockets API (“connect”, “send”, “recv”, “gethostbyname”
  etc). Network traffic from the VM looks like regular Mac application traffic
  and should get forwarded correctly by VPN software. We also support a “native”
  networking mode which uses the in-kernel facilities (via vmnet.framework)
  which is simpler but is more vulnerable to VPN software interference.
- a transparent docker proxy that intercepts API calls, rewrites requests and
  responses and performs side-effects to authorize volume mounts and manage
  port forwards.

Talking to docker
=================  

When the user types “docker ps” on the host the control flow looks like this:

![Control flow](diagrams/control.png)

- the user types “docker <command>”
- the docker CLI connect()s to /var/run/docker.sock (/var/run/docker.sock is a
  symlink to a per-user docker socket)
- the *docker proxy* goroutine inside the driver process
  (“com.docker.driver.amd64-linux”) accept()s the connection and reads the request
  - if the request is a start containing Binds, a 9P connection may be made to
    the control interface of the file server to authorise the request
  - if the request is a start containing Ports, a 9P connection may be made to
    the control interface of the port forwarder to bind() and listen()
- the *docker proxy* goroutine connect()s to an internal socket (the “underlying” socket)
- the docker proxy goroutine write()s the request to the underlying socket
- this is proxied by the hyperkit virtio-vsock backend into a connection to vsock
  port 2376 in the VM.
- the *vsudd* process in the Linux VM receives the connection
- the *vsudd* writes the request to the real docker daemon.

The Docker proxy
================

The docker proxy currently runs on the host but we may move some parts of it
into the Linux VM. The proxy intercepts the following API calls:
- start container: the proxy examines the Binds and Ports and
  - for each Bind, the paths are rewritten according to the following rules:
    - if the path is “/var/run/docker.sock” and the experimental
      expose-docker-socket configuration variable is set, the path is
      untransformed, effectively granting access to the container to the raw
      docker socket in the VM. This is obviously unsatisfactory, since we would
      like to expose a more limited interface, however it is in widespread use.
    - if the path begins with “/”, the proxy calls a 9P control interface
      exposed by the fileserver process “com.docker.osxfs” to authorize the
      mount and register for FSevents. The path is prefixed with “/Mac” since
      this is the mountpoint inside the VM
    - any other path is assumed to refer to a volume and is untransformed
  - for each Port, the ports are rewritten according to the following rules:
    - if the host port (and optionally IP) is specified and the experimental
      native-port-forwarding configuration variable is set, the proxy calls a 9P
      control interface exposed by the port forwarding process
      “com.docker.slirp” to call “bind”, “listen”, “accept” etc. The port
      exposed on the host is the same as the one exposed internally on the Linux VM
    - if the host port is not specified (meaning “choose one for me”), this is
      left untransformed introspect container: the proxy examines the results
      from Docker and reverses the path and port transforms.
- list containers: (slightly hacky!) if the native-port-forwarding is enabled
  and a host port was unspecified (meaning “choose one for me”) then the proxy
  contacts the 9P control interface of the port forwarding service and asks it
  to bind a fresh port on the Mac. This port is then substituted for the internal Linux port. We should probably make this a side-effect of start
  container instead.

The proxy listens for events and cleans up Binds and Port state when it
receives a container "die" event.

