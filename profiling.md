# Profiling Orca

If you run the Orca server with the debug flag set, not only will you get more logging output, but we enable
remote pprof access.

Links:
* http://blog.golang.org/profiling-go-programs
* https://golang.org/pkg/net/http/pprof/


# Examples

* First deploy orca with debug.
* If you're using TLS (e.g., with bootstrap install) you'll need to add the certs to the local system's trusted certs (unfortunately pprof doesn't have an "--insecure" flag)
    ```bash
sudo bash -c "docker run --rm -it \
        --name orca-bootstrap \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dockerorca/orca-bootstrap \
        dump-certs > /usr/local/share/ca-certificates/orca.crt"
sudo update-ca-certificates
```

* Now you can point pprof at the server

```bash
ORCA=https://192.168.104.73:443
go tool pprof ${ORCA}/debug/pprof/block

web
```

That will pop up a nice SVG image on your default browser showing the accumulated blocking calls

* Other URL endpoints of interest:  ("web" should produce a nice summary in each of these)
    * go tool pprof ${ORCA}/debug/pprof/profile     - The CPU usage
    * go tool pprof ${ORCA}/debug/pprof/heap     - The memory usage
    * go tool pprof ${ORCA}/debug/pprof/goroutine     - The goroutine usage
    * curl --insecure ${ORCA}/debug/pprof/   - Display the entry points (or use your browser)

# Dropped nodes

* It looks like nodefraction can be used to include nodes that would otherwise be dropped, but the web visualization seems to ignore it
    ```bash
go tool pprof -nodefraction=0.00000001 ${ORCA}/debug/pprof/block

peek swarm
```

