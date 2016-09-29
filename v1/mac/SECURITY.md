Docker for Mac security notes
=============================

Running docker commands
=======================

On Linux, only members of the `docker` group can usually execute docker
commands, and the `docker` group is considered equivalent to root.

On Docker for Mac, users who can execute docker commands have the following
extra capabilities:

Binding privileged ports
------------------------

When running in VPN proxy mode, it is possible to bind privileged ports
on the host (via `docker run -p 80:80 nginx`) via the privileged helper
process `com.docker.vmnetd`.

* Any user who has access to the `com.docker.vmnetd` socket can bind
privileged ports.

