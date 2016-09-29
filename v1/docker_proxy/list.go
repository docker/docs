package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/docker/engine-api/types"
)

type psRewriter struct {
	portRewriter PortRewriter
}

// NewPsRewriter returns a new Rewriter
func NewPsRewriter(portRewriter PortRewriter) Rewriter {
	if portRewriter == nil {
		return &nopRewriter{}
	}

	return &psRewriter{
		portRewriter: portRewriter,
	}
}

func (p *psRewriter) Rewrite(body io.ReadCloser) (int, io.ReadCloser) {
	var cs []types.Container
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&cs); err != nil {
		log.Printf("Failed to decode []types.Container: %#v\n", err)
		panic(err)
	}

	if p.portRewriter != nil {
		for i, container := range cs {
			cs[i].Ports = p.portRewriter.RewritePorts(container.ID, container.Ports)
		}
	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(&cs); err != nil {
		log.Printf("Failed to re-encode json: %#v\n", err)
		panic(err)
	}

	return b.Len(), ioutil.NopCloser(&b)
}
