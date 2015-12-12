+++
title = "Profiling UCP"
description = "Docker Universal Control Plane"
[menu.ucp]
weight="1"
+++


# Profiling UCP

If you run the UCP server with the debug flag set, not only will you get more
logging output, but we enable remote
[`pprof`](https://golang.org/pkg/net/http/pprof/) access. Because UCP is a go
program, it is a good idea to make yourself familiar with [profiling go
programs](http://blog.golang.org/profiling-go-programs).


# Examples

First deploy UCP with debug. If you're using TLS (e.g., with bootstrap install) you'll need to add the certs to the local system's trusted certs (unfortunately `pprof` doesn't have an `--insecure` flag)

```bash
sudo bash -c "docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dockerorca/ucp \
        dump-certs > /usr/local/share/ca-certificates/orca.crt"
sudo update-ca-certificates
```

Now you can point `pprof` at the server

```bash
UCP=https://192.168.104.73:443
go tool pprof ${UCP}/debug/pprof/block

web
```

That will pop up a nice SVG image on your default browser showing the accumulated blocking calls

Other URL endpoints of interest:  (`web` should produce a nice summary in each of these)

    * `go tool pprof ${UCP}/debug/pprof/profile`     - The CPU usage
    * `go tool pprof ${UCP}/debug/pprof/heap`     - The memory usage
    * `go tool pprof ${UCP}/debug/pprof/goroutine`     - The goroutine usage
    * `curl --insecure ${UCP}/debug/pprof/`   - Display the entry points (or use your browser)

# Dropped nodes

It looks like `nodefraction` can be used to include nodes that would otherwise be dropped, but the web visualization seems to ignore it

```bash
go tool pprof -nodefraction=0.00000001 ${UCP}/debug/pprof/block

peek swarm
```
