package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var dockerSocket string = "/var/run/docker.sock"

// ProxyServer is invoked by the `proxy` verb
func ProxyServer(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// Check if the docker socket is actually bound to the container
	dockerSocket = c.String("d")
	if _, err := os.Stat(dockerSocket); os.IsNotExist(err) {
		log.Fatal("Unable to locate the Docker Socket at /var/run/docker.sock")
	}

	// Check if TLS settings have been specified
	if c.String("cert") == "" || c.String("key") == "" {
		log.Fatal("Unable to start ucp-proxy without TLS configuration")
	}

	// Create the http.Server and establish a redirect of all routes to the engine
	http.HandleFunc("/", engineRedirect)
	server := &http.Server{
		Addr: c.String("listen-address"),
	}

	// Configure mTLS
	log.Infof("Configuring TLS: ca=%s cert=%s key=%s", c.String("ca"), c.String("cert"), c.String("key"))
	caCert, err := ioutil.ReadFile(c.String("ca"))
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	server.TLSConfig = &tls.Config{
		// ListenAndServeTLS will wire up the cert/key pairs automatically
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}

	// Block forever until an error occurs
	return server.ListenAndServeTLS(c.String("cert"), c.String("key"))
}

func engineRedirect(w http.ResponseWriter, r *http.Request) {
	// TODO: increment an active request counter
	// TODO: reject new requests in shutdown-mode
	var c net.Conn

	cl, err := net.Dial("unix", dockerSocket)
	if err != nil {
		log.Errorf("error connecting to backend: %s", err)
		return
	}

	c = cl
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijack error", 500)
		return
	}
	nc, _, err := hj.Hijack()
	if err != nil {
		log.Printf("hijack error: %v", err)
		return
	}
	defer nc.Close()
	defer c.Close()

	err = r.Write(c)
	if err != nil {
		log.Printf("error copying request to target: %v", err)
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go cp(c, nc)
	go cp(nc, c)
	<-errc
	//TODO: decrement active request counter
}
