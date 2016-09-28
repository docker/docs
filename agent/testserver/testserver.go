package testserver

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// TestServer is invoked for the `test-server` command
func TestServer(c *cli.Context) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	log.Infof("Listening on %s", c.String("listen-address"))
	log.Fatal(http.ListenAndServe(c.String("listen-address"), nil))
	return nil
}
