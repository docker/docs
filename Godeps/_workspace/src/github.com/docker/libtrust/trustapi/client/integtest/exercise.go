package main

import (
	"log"

	"github.com/docker/libtrust"
	"github.com/docker/libtrust/trustapi/client"
)

func exerciseBaseGraph(c *client.TrustClient, name string) []byte {
	statement, err := c.GetBaseGraph(name)
	if err != nil {
		log.Fatal(err)
	}

	b, err := statement.Bytes()
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func main() {
	pool, err := libtrust.LoadCertificatePool("./ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	c := client.NewTrustClient("localhost:8092", pool)
	b1 := exerciseBaseGraph(c, "empty")
	b2 := exerciseBaseGraph(c, "dmcgowan")
	_, err = c.GetBaseGraph("empty-bad")
	if err == nil {
		log.Fatalf("Did not receive error getting empty-bad")
	}
	log.Printf("Expected error getting empty-bad: %s", err)

	log.Printf("Statements:\nempty:\n%s\n\ndmcgowan\n%s", b1, b2)
}
