package nio_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/oliverpool/nio"
)

func TestPipeClose(t *testing.T) {
	var buf bytes.Buffer
	p := nio.NewPipe(&buf)
	a := errors.New("a")
	b := errors.New("b")
	p.CloseWithError(a)
	p.CloseWithError(b)
	_, err := p.Read(make([]byte, 1))
	if err != a {
		t.Errorf("err = %v want %v", err, a)
	}
}

func TestBigWriteSmallBuf(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 5))
	p := nio.NewPipe(buf)

	go func() {
		defer p.Close()
		n, err := p.Write([]byte("hello world"))
		if err != nil {
			t.Error(err)
		}
		if int(n) != len("hello world") {
			t.Errorf("wrote wrong # of bytes")
		}
	}()

	output := bytes.NewBuffer(nil)
	_, err := io.Copy(output, p)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(output.Bytes(), []byte("hello world")) {
		t.Errorf("unexpected output %s", output.Bytes())
	}
}

func TestPipeCloseEarly(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	p := nio.NewPipe(buf)
	p.Close()

	_, err := p.Write([]byte("hello world"))
	if err != io.ErrClosedPipe {
		t.Errorf("expected closed pipe, got %v", err)
	}

	_, err = io.Copy(ioutil.Discard, p)
	if err != nil {
		t.Errorf("expected nil, got %v")
	}
}

func TestPipe(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	p := nio.NewPipe(buf)

	data := []byte("the quick brown fox jumps over the lazy dog")
	if _, err := p.Write(data); err != nil {
		t.Error(err)
		return
	}
	p.Close()

	result := make([]byte, 1024)
	n, err := p.Read(result)
	if err != nil {
		t.Error(err)
		return
	}
	result = result[:n]

	if !bytes.Equal(data, result) {
		t.Errorf("exp [%s]\ngot[%s]", string(data), string(result))
	}
	if n, err := p.Read(result); err != io.EOF || n != 0 {
		t.Error(n, err)
		return
	}
}

func TestEarlyCloseWrite(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1))
	p := nio.NewPipe(buf)

	testerr := errors.New("test err")

	p.CloseWithError(testerr)

	_, err := p.Write([]byte(""))

	if err != io.ErrClosedPipe {
		t.Errorf("expected %s but got %s.", io.ErrClosedPipe, err)
	}

	_, err = io.Copy(ioutil.Discard, p)
	if err != testerr {
		t.Errorf("expected %s but got %s.", testerr, err)
	}
}

type badBuffer struct{}

func (badBuffer) Len() int                    { return 3 }
func (badBuffer) Write(p []byte) (int, error) { return len(p), nil }
func (badBuffer) Read(p []byte) (int, error)  { return 0, io.EOF }

func TestEmpty(t *testing.T) {
	p := nio.NewPipe(badBuffer{})
	n, err := p.Write([]byte("any"))

	if err != nil {
		t.Error(err)
	}

	if n != 3 {
		t.Errorf("wrote wrong # of bytes %d", n)
	}

	n, err = p.Read(nil)

	if err != io.EOF {
		t.Error(err)
	}

	if n != 0 {
		t.Errorf("wrote wrong # of bytes %d", n)
	}
}

func BenchmarkPipe(b *testing.B) {
	p := nio.NewPipe(bytes.NewBuffer(make([]byte, 0, 1024)))
	r := slowReader{p, false}
	benchPipe(r, p, b)
	r.fast = true
}

func BenchmarkIOPipe(b *testing.B) {
	pr, w := io.Pipe()
	r := slowReader{pr, false}
	benchPipe(r, w, b)
	r.fast = true
}

type slowReader struct {
	r    io.Reader
	fast bool
}

func (r slowReader) Read(data []byte) (int, error) {
	if !r.fast {
		time.Sleep(1000 * time.Nanosecond)
	}
	return r.r.Read(data)
}

func benchPipe(r io.Reader, w io.WriteCloser, b *testing.B) {
	b.ReportAllocs()
	go io.Copy(ioutil.Discard, r)

	for i := 0; i < b.N; i++ {
		w.Write([]byte("hello world"))
	}
	w.Close()
}
