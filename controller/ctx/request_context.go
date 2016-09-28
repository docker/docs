package ctx

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/resources"
)

var ErrHTTPRequestNotSet = errors.New("the http.Request has not been set for the provided OrcaRequestContext")

type (
	OrcaRequestContext struct {
		// Auth contains the authenticated account used to authorize
		// the incoming request
		Auth *auth.Context

		Request *http.Request

		// PathVars contains the path variables of the request
		PathVars map[string]string
		// QueryVars contains the query variables of the request
		QueryVars url.Values

		// RequiresNotary is set to true if the request requires verifying the image signature with notary
		RequiresNotary bool

		// bodyBuffer contains the Body of the request, if it has been read
		bodyBuffer []byte

		// readBody is set to true if Request.Body has been read
		readBody bool

		// Resources is the set of resources that are being requested.
		// Resources are populated by request parsers and are used by the access control
		// middleware layer
		Resources []resources.ResourceRequest

		// MainResource is a reference to the primary resource of a request.
		// This is meant to differentiate encapsulated resources from the top-level resource
		// and to also allow easier extraction of concrete resources in handlers
		MainResource resources.ResourceRequest
	}
)

// The parseBody method populates rc.bodyBuffer with the contents of the request
// and then proceeds to re-populate the request Body, in case it is forwarded
func (rc *OrcaRequestContext) parseBody() error {
	var err error

	if rc.Request == nil {
		return ErrHTTPRequestNotSet
	}
	if rc.Request.Body == nil {
		rc.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
	}

	rc.bodyBuffer, err = ioutil.ReadAll(rc.Request.Body)
	// Do not defer Body.Close as we will repopulate the body
	rc.Request.Body.Close()
	if err != nil {
		return err
	}

	rc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rc.bodyBuffer))
	rc.Request.ContentLength = int64(len(rc.bodyBuffer))
	rc.readBody = true
	return nil
}

// ParseVars populates rc.PathVars and rc.QueryVars from the HTTP Request
func (rc *OrcaRequestContext) ParseVars() error {
	if rc.Request == nil {
		return ErrHTTPRequestNotSet
	}
	rc.PathVars = mux.Vars(rc.Request)
	if rc.PathVars == nil {
		rc.PathVars = make(map[string]string)
	}

	rc.QueryVars = rc.Request.URL.Query()
	return nil
}

// Body() generates an io.Reader from the parsed rc.bodyBuffer
func (rc *OrcaRequestContext) Body() io.Reader {
	if !rc.readBody {
		err := rc.parseBody()
		if err != nil {
			log.Error(err)
		}
	}
	return bytes.NewReader(rc.bodyBuffer)
}

// BodyBuffer() generates a []byte from the parsed rc.bodyBuffer
func (rc *OrcaRequestContext) BodyBuffer() []byte {
	if !rc.readBody {
		err := rc.parseBody()
		if err != nil {
			log.Error(err)
		}
	}
	return rc.bodyBuffer
}
