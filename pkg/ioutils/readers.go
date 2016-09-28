package ioutils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"
	"sync"
	"time"
)

type readCloserWrapper struct {
	io.Reader
	closer func() error
}

func (r *readCloserWrapper) Close() error {
	return r.closer()
}

// NewReadCloserWrapper returns a new io.ReadCloser.
func NewReadCloserWrapper(r io.Reader, closer func() error) io.ReadCloser {
	return &readCloserWrapper{
		Reader: r,
		closer: closer,
	}
}

type readerErrWrapper struct {
	reader io.Reader
	closer func()
}

func (r *readerErrWrapper) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		r.closer()
	}
	return n, err
}

// NewReaderErrWrapper returns a new io.Reader.
func NewReaderErrWrapper(r io.Reader, closer func()) io.Reader {
	return &readerErrWrapper{
		reader: r,
		closer: closer,
	}
}

// bufReader allows the underlying reader to continue to produce
// output by pre-emptively reading from the wrapped reader.
// This is achieved by buffering this data in bufReader's
// expanding buffer.
type bufReader struct {
	sync.Mutex
	buf      io.ReadWriter
	reader   io.Reader
	err      error
	wait     sync.Cond
	drainBuf []byte
}

// NewBufReader returns a new bufReader.
func NewBufReader(r io.Reader) io.ReadCloser {
	reader := &bufReader{
		buf:      NewBytesPipe(nil),
		reader:   r,
		drainBuf: make([]byte, 1024),
	}
	reader.wait.L = &reader.Mutex
	go reader.drain()
	return reader
}

// NewBufReaderWithDrainbufAndBuffer returns a BufReader with drainBuffer and buffer.
func NewBufReaderWithDrainbufAndBuffer(r io.Reader, drainBuffer []byte, buffer io.ReadWriter) io.ReadCloser {
	reader := &bufReader{
		buf:      buffer,
		drainBuf: drainBuffer,
		reader:   r,
	}
	reader.wait.L = &reader.Mutex
	go reader.drain()
	return reader
}

func (r *bufReader) drain() {
	var (
		duration       time.Duration
		lastReset      time.Time
		now            time.Time
		reset          bool
		bufLen         int64
		dataSinceReset int64
		maxBufLen      int64
		reuseBufLen    int64
		reuseCount     int64
	)
	reuseBufLen = int64(len(r.reuseBuf))
	lastReset = time.Now()
	for {
		//Call to scheduler is made to yield from this goroutine.
		//This avoids goroutine looping here when n=0,err=nil, fixes code hangs when run with GCC Go.
		callSchedulerIfNecessary()
		n, err := r.reader.Read(r.drainBuf)
		dataSinceReset += int64(n)
		r.Lock()
		bufLen = int64(r.buf.Len())
		if bufLen > maxBufLen {
			maxBufLen = bufLen
		}

		// Avoid unbounded growth of the buffer over time.
		// This has been discovered to be the only non-intrusive
		// solution to the unbounded growth of the buffer.
		// Alternative solutions such as compression, multiple
		// buffers, channels and other similar pieces of code
		// were reducing throughput, overall Docker performance
		// or simply crashed Docker.
		// This solution releases the buffer when specific
		// conditions are met to avoid the continuous resizing
		// of the buffer for long lived containers.
		//
		// Move data to the front of the buffer if it's
		// smaller than what reuseBuf can store
		if bufLen > 0 && reuseBufLen >= bufLen {
			n, _ := r.buf.Read(r.reuseBuf)
			r.buf.Write(r.reuseBuf[0:n])
			// Take action if the buffer has been reused too many
			// times and if there's data in the buffer.
			// The timeout is also used as means to avoid doing
			// these operations more often or less often than
			// required.
			// The various conditions try to detect heavy activity
			// in the buffer which might be indicators of heavy
			// growth of the buffer.
		} else if reuseCount >= r.maxReuse && bufLen > 0 {
			now = time.Now()
			duration = now.Sub(lastReset)
			timeoutReached := duration >= r.resetTimeout

			// The timeout has been reached and the
			// buffered data couldn't be moved to the front
			// of the buffer, so the buffer gets reset.
			if timeoutReached && bufLen > reuseBufLen {
				reset = true
			}
			// The amount of buffered data is too high now,
			// reset the buffer.
			if timeoutReached && maxBufLen >= r.bufLenResetThreshold {
				reset = true
			}
			// Reset the buffer if a certain amount of
			// data has gone through the buffer since the
			// last reset.
			if timeoutReached && dataSinceReset >= r.maxReadDataReset {
				reset = true
			}
			// The buffered data is moved to a fresh buffer,
			// swap the old buffer with the new one and
			// reset all counters.
			if reset {
				newbuf := &bytes.Buffer{}
				newbuf.ReadFrom(r.buf)
				r.buf = newbuf
				lastReset = now
				reset = false
				dataSinceReset = 0
				maxBufLen = 0
				reuseCount = 0
			}
		}
		if err != nil {
			r.err = err
		} else {
			if n == 0 {
				// nothing written, no need to signal
				r.Unlock()
				continue
			}
			r.buf.Write(r.drainBuf[:n])
		}
		reuseCount++
		r.wait.Signal()
		r.Unlock()
		if err != nil {
			break
		}
	}
}

func (r *bufReader) Read(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()
	for {
		n, err = r.buf.Read(p)
		if n > 0 {
			return n, err
		}
		if r.err != nil {
			return 0, r.err
		}
		r.wait.Wait()
	}
}

// Close closes the bufReader
func (r *bufReader) Close() error {
	closer, ok := r.reader.(io.ReadCloser)
	if !ok {
		return nil
	}
	return closer.Close()
}

// HashData returns the sha256 sum of src.
func HashData(src io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

// OnEOFReader wraps a io.ReadCloser and a function
// the function will run at the end of file or close the file.
type OnEOFReader struct {
	Rc io.ReadCloser
	Fn func()
}

func (r *OnEOFReader) Read(p []byte) (n int, err error) {
	n, err = r.Rc.Read(p)
	if err == io.EOF {
		r.runFunc()
	}
	return
}

// Close closes the file and run the function.
func (r *OnEOFReader) Close() error {
	err := r.Rc.Close()
	r.runFunc()
	return err
}

func (r *OnEOFReader) runFunc() {
	if fn := r.Fn; fn != nil {
		fn()
		r.Fn = nil
	}
}
