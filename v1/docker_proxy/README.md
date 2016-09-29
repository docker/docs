# Proxy server for Docker

This is a small server that listens on one socket (typically `/var/run/docker.sock` on OS X) and forwards requests to another socket (e.g. one that leads to the real Docker running under XHyve).

While doing this, it modifies host volume paths from `/path` to `/Mac/path` so that they are correct for the XHyve environment.

TODO: It should also notify the file sharing service about which paths need to be shared and for which containers, allowing the sharing service to prevent access to files that have not been shared. This is a protection against a compromised container getting access to the boot2docker root and from there to the whole OS X filesystem. For this reason, the proxy server cannot itself run under XHyve.

## Testing

To test on Linux, forwarding to the real Docker daemon:

    $ go run proxy.go -verbose \
        -listen /tmp/proxy.sock \
        -underlying /var/run/docker.sock \
    INFO[0000] Listening on /tmp/proxy.sock

Run the docker client:

    $ docker -H unix:///tmp/proxy.sock run -it debian:jessie /bin/bash
    #


## See also

The Docker repository includes support for writing middleware that intercepts Docker API calls. Currently it appears that this only works when run in the same process as Docker itself, but in future it may be useful to build on that code.
