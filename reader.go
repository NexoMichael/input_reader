// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

package inputreader

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

	newState := *termios

	newState.Lflag |= unix.ISIG
	newState.Lflag &^= (syscall.ICANON | syscall.IEXTEN)

	newState.Iflag |= unix.ICRNL
	newState.Iflag &^= (syscall.PARMRK | syscall.ISTRIP | syscall.IXON)

	newState.Oflag &^= syscall.OPOST

	newState.Cflag |= syscall.CS8
	newState.Cflag &^= (syscall.CSIZE | syscall.PARENB)

	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		return nil, err
	}

	return &InputReader{fd: fd, oldState: termios}, nil
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
