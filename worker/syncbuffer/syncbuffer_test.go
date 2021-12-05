package syncbuffer

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBufferSimple(t *testing.T) {
	content := "Hello, world!"
	b := New()
	r := b.NewReader()
	if n, err := b.Write([]byte(content)); n != len([]byte(content)) || err != nil {
		t.Errorf("Buffer.Write() = %d, %v; want %d, nil", n, err, len([]byte(content)))
	}
	if err := b.Close(); err != nil {
		t.Errorf("Buffer.Close() returned non-nil error: %v", err)
	}
	if d, err := ioutil.ReadAll(r); string(d) != content || err != nil {
		t.Errorf("io.ReadAll() = %q, %v; want %q, nil", string(d), content, err)
	}
	if n, err := r.Read(make([]byte, 1)); n != 0 || err != io.EOF {
		t.Errorf("After EOF, Read() = %d, %v; want 0, EOF", n, err)
	}
}

func TestBufferInterleavingReadAndWrite(t *testing.T) {
	content := []string{"Hello, ", "world!"}
	b := New()
	r := b.NewReader()
	for i, c := range content {
		wantN := len([]byte(c))
		if n, err := b.Write([]byte(c)); n != wantN || err != nil {
			t.Errorf("[%d] Buffer.Write(%q) = %d, %v; want %d, nil", i, c, n, err, wantN)
		}
		buf := make([]byte, 100)
		if n, err := r.Read(buf); n != wantN || err != nil {
			t.Errorf("[%d] Read() = %d, %v; want %d, nil", i, n, err, wantN)
		}
	}
	if err := b.Close(); err != nil {
		t.Errorf("Buffer.Close() returned non-nil error: %v", err)
	}
	if n, err := r.Read(make([]byte, 1)); n != 0 || err != io.EOF {
		t.Errorf("After EOF, Read() = %d, %v; want 0, EOF", n, err)
	}
	// A new reader should be able to read from the beginning.
	r2 := b.NewReader()
	want := strings.Join(content, "")
	if got, err := ioutil.ReadAll(r2); string(got) != want || err != nil {
		t.Errorf("io.ReadAll() = %q, %v; want %q, nil", string(want), content, err)
	}
}

func TestBufferInParallel(t *testing.T) {
	a := []byte("a")
	b := New()
	// Add a constant to make sure we have some readers/writers.
	nReaders := rand.Intn(1000) + 10
	nWrites := rand.Intn(1000) + 10
	wantContent := make([]byte, 0, len(a)*nReaders)
	for i := 0; i < nWrites; i++ {
		wantContent = append(wantContent, a...)
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func(w io.WriteCloser) {
		defer wg.Done()
		for i := 0; i < nWrites; i++ {
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))
			if n, err := w.Write(a); n != len(a) || err != nil {
				t.Errorf("[%d] Write() = %d, %v; want %d, nil", i, n, err, len(a))
			}
		}
		if err := w.Close(); err != nil {
			t.Errorf("Close() returned non-nil error: %v", err)
		}
	}(b)
	for i := 0; i < nReaders; i++ {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(10)))
		wg.Add(1)
		go func(i int, r io.Reader) {
			defer wg.Done()
			gotContent, err := io.ReadAll(r)
			if err != nil {
				t.Errorf("[Reader %d] ReadAll() returned non-nil error: %v", i, err)
			}
			if bytes.Compare(gotContent, wantContent) != 0 {
				t.Errorf("[Reader %d] ReadAll() got unexpected content", i)
			}
		}(i, b.NewReader())
	}
	wg.Wait()
}
