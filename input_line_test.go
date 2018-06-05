package input_reader

import (
	"errors"
	"strings"
	"testing"
)

func TestDefaultBuffer(t *testing.T) {
	buf := NewBuffer()
	if buf == nil {
		t.Error("buffer is not created")
	}

	if 4096 != cap(buf.buf) {
		t.Error("buffer default capacity should be 4KB")
	}
	if 0 != len(buf.buf) {
		t.Error("buffer default length should be 0")
	}
}

func TestGetBufferConcurrent(t *testing.T) {
	k := 1000000
	// update loop
	buf := NewBuffer()
	go func() {
		for i := 0; i < k; i++ {
			buf.Lock()
			buf.buf = append(buf.buf, 'a')
			if i%1000 == 0 {
				buf.buf = buf.buf[:0]
			}
			buf.Unlock()
		}
	}()
	// read loop
	go func() {
		for i := 0; i < k; i++ {
			if len(buf.Buffer()) > 1000 {
				t.Error("something goes wrong")
			}
		}
	}()
}

func TestReadLine(t *testing.T) {
	buf := NewBuffer()
	str := "some string\r\n"
	r := strings.NewReader(str)
	res, err := buf.ReadLine(r)
	if err != nil {
		t.Error("failed to read line", "err", err)
	}
	if string(res) != str[:len(str)-2] {
		t.Error("readed different string", "expected", str, "actual", string(res))
	}
}

type badReader struct{}

func (r *badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some error")
}
func TestBadReader(t *testing.T) {
	buf := NewBuffer()
	_, err := buf.ReadLine(&badReader{})
	if err == nil {
		t.Error("bad reader allowed to read line")
	}
}

type closedReader struct{}

func (r *closedReader) Read(p []byte) (n int, err error) {
	return 0, nil
}
func TestClosedReader(t *testing.T) {
	buf := NewBuffer()
	res, err := buf.ReadLine(&closedReader{})
	if err != nil {
		t.Error("failed to read line", "err", err)
	}
	if len(res) != 0 {
		t.Error("bytes should not be readed")
	}
}
