package syncbuffer

import (
	"io"
	"sync"
)

// Buffer is a synchronized 1-writer-N-reader unbounded in-memory buffer.
//
// The Buffer itself is the writer. Use NewReader to create a reader. There is
// only one copy of the internal data. Synchronization details are hidden
// behind Read, Write and Close, all of which may block.
//
// Do not copy this type. Do not use the zero value; use New() instead.
type Buffer struct {
	buf []byte
	eof bool

	wlock sync.RWMutex
	rcond *sync.Cond
}

// New creates a new Buffer.
func New() *Buffer {
	buffer := new(Buffer)
	buffer.rcond = sync.NewCond(buffer.wlock.RLocker())
	return buffer
}

// Write implements io.Writer. Write will block until the Buffer is writable.
// Write will always succeed unless OOM (in which case the program panics).
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.wlock.Lock()
	defer b.wlock.Unlock()
	if b.eof {
		// Nothing is changed here, so we don't need to broadcast.
		return 0, io.ErrClosedPipe
	}
	b.buf = append(b.buf, p...)
	b.rcond.Broadcast()
	return len(p), nil
}

// Close implements io.Closer. Close will block until the Buffer is writable.
// Close will always succeed.
func (b *Buffer) Close() error {
	b.wlock.Lock()
	defer b.wlock.Unlock()
	b.eof = true
	b.rcond.Broadcast()
	return nil
}

// NewReader creates a new io.Reader with an independent seek position.
func (b *Buffer) NewReader() io.Reader {
	return &syncReader{b: b, c: b.rcond}
}

type syncReader struct {
	b   *Buffer
	pos int
	c   *sync.Cond
}

func (r *syncReader) Read(p []byte) (n int, err error) {
	r.c.L.Lock()
	// Condition: there is more in the buffer to read or we've reached EOF.
	for !(r.pos < len(r.b.buf) || r.b.eof) {
		r.c.Wait()
	}
	defer r.c.L.Unlock()
	n = copy(p, r.b.buf[r.pos:])
	r.pos += n
	// Only return EOF when we've read everything and there won't be more.
	if r.pos == len(r.b.buf) && r.b.eof {
		err = io.EOF
	}
	return n, err
}
