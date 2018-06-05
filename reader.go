// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

package inputreader // import "github.com/NexoMichael/inputreader"

import (
	"errors"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// InputReader provides ReaderCloser interface for standart input or other input source pipe.
type InputReader struct {
	oldState *unix.Termios
	fd       int
	stop     bool
}

// NewInputReader returns a new InputReader for *os.File pipe.
func NewInputReader(in *os.File) (*InputReader, error) {
	if in == nil {
		return nil, errors.New("input is not set")
	}
	fd := int(in.Fd())

	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	newState := newTermios(*termios)

	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		return nil, err
	}

	return &InputReader{fd: fd, oldState: termios}, nil
}

func newTermios(termios unix.Termios) unix.Termios {
	termios.Lflag &^= (unix.ISIG | syscall.ICANON | syscall.IEXTEN)

	termios.Iflag |= unix.ICRNL
	termios.Iflag &^= (syscall.PARMRK | syscall.ISTRIP | syscall.IXON)

	termios.Oflag &^= syscall.OPOST

	termios.Cflag |= syscall.CS8
	termios.Cflag &^= (syscall.CSIZE | syscall.PARENB)

	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	return termios
}

// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence n may be less than len(p).
// At EOF, the count will be zero and err will be io.EOF.
func (r *InputReader) Read(p []byte) (n int, err error) {
	return syscall.Read(int(r.fd), p[:])
}

// Close resets terminal settings to the initial state.
func (r *InputReader) Close() (err error) {
	unix.IoctlSetTermios(r.fd, ioctlWriteTermios, r.oldState)
	return err
}
