package nio

import (
	"io"
	"sync"
)

// Pipe holds a buffered pipe.
type Pipe struct {
	b   Buffer
	c   sync.Cond
	m   sync.Mutex
	err error
}

// Buffer is the underlying buffer for a Pipe.
// Calls to Read, Write and Len will be sequential (no need for thread safety).
type Buffer interface {
	io.ReadWriter
	Len() int
}

// NewPipe creates a new Pipe (safe for concurrency) based on the given buffer.
//
// It can be used to connect code expecting an io.Reader with code expecting an io.Writer.
// Reads on one end read from the supplied Buffer. Writes write to the supplied Buffer.
// It is safe to call Read and Write in parallel or with Close.
// Close will complete once pending I/O is done, and may cancel blocking Read/Writes.
// Buffered data will still be available to Read after the Writer has been closed.
// Parallel calls to Read, and parallel calls to Write are not safe.
func NewPipe(b Buffer) *Pipe {
	var p Pipe
	p.b = b
	p.c.L = &p.m
	return &p
}

// Read from the buffer into b.
// Blocks if the buffer is empty (and not closed)
func (p *Pipe) Read(b []byte) (n int, err error) {
	p.c.L.Lock()
	defer p.c.L.Unlock()
	for p.b.Len() == 0 {
		if p.err != nil {
			return 0, p.err
		}
		p.c.Wait()
	}
	n, err = p.b.Read(b)
	return
}

// Write copies bytes from b into the buffer and wakes a reader.
func (p *Pipe) Write(b []byte) (n int, err error) {
	p.c.L.Lock()
	defer p.c.L.Unlock()
	if p.err != nil {
		return 0, io.ErrClosedPipe
	}
	defer p.c.Signal()
	return p.b.Write(b)
}

// Close closes the pipe (reads will succeed until exhaustion)
func (p *Pipe) Close() error {
	return p.CloseWithError(nil)
}

// CloseWithError closes the pipe
// (which will be returned on Read happening after the exhaustion)
func (p *Pipe) CloseWithError(err error) error {
	if err == nil {
		err = io.EOF
	}

	p.c.L.Lock()
	defer p.c.L.Unlock()
	if p.err != nil {
		return io.ErrClosedPipe
	}

	defer p.c.Signal()
	p.err = err
	return nil
}
