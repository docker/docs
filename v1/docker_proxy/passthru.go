package proxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Rewriter allows for an API request or response to be rewritten
type Rewriter interface {
	Rewrite(io.ReadCloser) (int, io.ReadCloser)
}

type passthru struct {
	backendDialer BackendDialer
	rewriter      Rewriter
	cancellable   bool
}

func newPassthru(backendDialer BackendDialer, rewriter Rewriter) *passthru {
	return &passthru{
		backendDialer: backendDialer,
		rewriter:      rewriter,
		cancellable:   false,
	}
}

func newCancellablePassthru(backendDialer BackendDialer, rewriter Rewriter) *passthru {
	return &passthru{
		backendDialer: backendDialer,
		rewriter:      rewriter,
		cancellable:   true,
	}
}

// WriteFlusher extends the io.Writer interface with http.Flusher
type WriteFlusher interface {
	io.Writer
	http.Flusher
}

type writeFlusher struct {
	writer WriteFlusher
}

func (t *writeFlusher) Write(buf []byte) (int, error) {
	count, err := t.writer.Write(buf)
	t.writer.Flush()
	return count, err
}

func isRawStreamUpgrade(resp *http.Response) bool {
	return (resp.StatusCode == 101 && resp.Header.Get("Upgrade") == "tcp") ||
		resp.Header.Get("Content-Type") == "application/vnd.docker.raw-stream"
}

func (p *passthru) HandleHTTP(writer http.ResponseWriter, req *http.Request) error {
	log.Printf("proxy >> %s %s\n", req.Method, req.URL)

	// Connect to underlying service
	var underlying ReadWriteWriteCloser
	underlying, err := p.backendDialer.Dial()
	if err != nil {
		return err
	}
	defer underlying.Close()

	if p.cancellable {
		if notifier, ok := writer.(http.CloseNotifier); ok {
			notify := notifier.CloseNotify()
			finished := make(chan struct{})
			defer close(finished)
			go func() {
				select {
				case <-notify:
					log.Println("Cancel connection...")
					underlying.CloseWrite()
				case <-finished:
				}
			}()
		}
	}

	return p.doHandleHTTP(underlying, writer, req)
}

func (p *passthru) doHandleHTTP(underlying ReadWriteWriteCloser, writer http.ResponseWriter, req *http.Request) error {
	underlying = &withLogging{label: "underlying", underlying: underlying}

	// Forward request to underlying
	requestErrors := make(chan error)
	go func() {
		requestErrors <- req.Write(underlying)
	}()

	// Read response
	resp, err := http.ReadResponse(bufio.NewReader(underlying), req)
	if err != nil {
		log.Println("error reading response from Docker: ", err)
		http.Error(writer, "Bad response from Docker engine", 502)
		return nil // ??
	}

	// Forward response to client
	isRaw := isRawStreamUpgrade(resp)
	copyHeaders(writer.Header(), resp.Header)
	if isRaw {
		// Stop Go adding a chunked encoding header
		writer.Header().Set("Transfer-Encoding", "identity")
		writer.WriteHeader(resp.StatusCode)
		writer.(http.Flusher).Flush()

		// Attach console - hijack connection
		err := upgradeToRaw(writer, underlying)
		if err != nil {
			return err
		}
	} else {
		// Regular HTTP response
		newContentLength, body := p.rewriter.Rewrite(resp.Body)
		if newContentLength >= 0 {
			writer.Header().Set(http.CanonicalHeaderKey("content-length"), fmt.Sprintf("%d", newContentLength))
		}
		writer.WriteHeader(resp.StatusCode)
		writer.(http.Flusher).Flush()

		_, err := io.Copy(&writeFlusher{writer: writer.(WriteFlusher)}, body)
		if err != nil {
			log.Println("error copying response body from Docker: ", err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.Println("error closing response body from Docker: ", err)
		}
	}

	// Wait for request thread to finish if it's still going
	err = <-requestErrors
	if err != nil {
		log.Printf("error forwarding client's request to Docker: %s", err)
	}
	log.Printf("proxy << %s %s\n", req.Method, req.URL)

	return nil
}

type nopRewriter struct{}

func (n nopRewriter) Rewrite(body io.ReadCloser) (int, io.ReadCloser) {
	return -1, body
}

// CloseWriter allows for a write to be closed
type CloseWriter interface {
	CloseWrite() error
}

func upgradeToRaw(writer http.ResponseWriter, underlying io.ReadWriteCloser) error {
	log.Println("Upgrading to raw stream")
	hj, ok := writer.(http.Hijacker)
	if !ok {
		panic("BUG: webserver doesn't support hijacking")
	}

	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return err
	}

	defer conn.Close()
	bufrw.Flush()
	done := make(chan bool)

	// Stream underlying -> conn
	go func() {
		var err error
		n := bufrw.Reader.Buffered()
		if n > 0 {
			buf := make([]byte, n)
			n, err := bufrw.Read(buf)
			if err != nil {
				panic(err)
			}
			_, err = conn.Write(buf[0:n])
		}
		if err != nil {
			log.Print("Error draining buffer: ", err)
		} else {
			_, err := io.Copy(conn, underlying)
			if err != nil {
				log.Print("Error forwarding raw stream from container: ", err)
			}
			err = conn.(CloseWriter).CloseWrite()
			if err != nil {
				log.Print("Error closing raw stream from container: ", err)
			}
		}
		done <- true
	}()

	// Stream underlying <- conn
	_, err = io.Copy(underlying, conn)
	if err != nil {
		log.Print("Error forwarding raw stream to container: ", err)
	}

	err = underlying.(CloseWriter).CloseWrite()
	if err != nil {
		log.Print("Error closing raw stream to container: ", err)
	}
	<-done

	return nil
}

// Make dst have the same contents as src
func copyHeaders(dst, src http.Header) {
	// http://stackoverflow.com/questions/13812121/how-to-clear-a-map-in-go
	for k := range dst {
		delete(dst, k)
	}
	for k, v := range src {
		dst[k] = v
	}
}

// ReadWriteWriteCloser extends io.ReadWriteCloser with a CloseWrite call
type ReadWriteWriteCloser interface {
	io.ReadWriteCloser
	CloseWrite() error
}

// A ReadWriteCloser that wraps underlying and logs all reads and writes.
type withLogging struct {
	label      string
	underlying ReadWriteWriteCloser
}

func (t *withLogging) Read(buf []byte) (int, error) {
	count, err := t.underlying.Read(buf)
	if verbose {
		log.Printf("proxy %s -> %#v\n", t.label, string(buf[:count]))
		if err != nil {
			log.Printf("proxy %s -> %s\n", t.label, err)
		}
	}
	return count, err
}

func (t *withLogging) Write(buf []byte) (int, error) {
	if verbose {
		log.Printf("proxy %s <- %#v\n", t.label, string(buf))
	}
	return t.underlying.Write(buf)
}

func (t *withLogging) Close() error {
	if verbose {
		log.Printf("proxy %s <- EOF\n", t.label)
	}
	return t.underlying.Close()
}

func (t *withLogging) CloseWrite() error {
	return t.underlying.CloseWrite()
}
