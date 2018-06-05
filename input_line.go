package inputreader

import (
	"io"
	"sync"
)

// InputLine holds current eneterd line of symbols
type InputLine struct {
	buf []byte

	lock sync.Mutex
}

// NewBuffer returns new InputLine buffer with default size equal to 4 KB.
func NewBuffer() *InputLine {
	return &InputLine{
		buf: make([]byte, 0, 4096),
	}
}

// Buffer returns current entered synmbols
func (l *InputLine) Buffer() []byte {
	l.lock.Lock()
	defer l.lock.Unlock()
	buf := make([]byte, len(l.buf))
	copy(buf, l.buf)
	return buf
}

// ReadLine tries to return a single line, not including the end-of-line bytes.
func (l *InputLine) ReadLine(r io.Reader) ([]byte, error) {
	var b [1]byte
	for {
		n, err := r.Read(b[:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			break
		}

		switch b[0] {
		case '\r':
			continue
		case '\n':
			l.lock.Lock()
			buf := make([]byte, len(l.buf))
			copy(buf, l.buf)
			l.buf = l.buf[:0]
			l.lock.Unlock()
			return buf, nil
		default:
			l.lock.Lock()
			l.buf = append(l.buf, b[0])
			l.lock.Unlock()
		}
	}
	return nil, nil
}
