package inputreader // import "github.com/NexoMichael/inputreader"

import (
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestBadInputReader(t *testing.T) {
	_, err := NewInputReader(nil)
	if err == nil {
		t.Error("input reader should not be created")
	}

	_, err = NewInputReader(os.Stdout)
	if err == nil {
		t.Error("input reader should not be created")
	}
}

func TestNewTermios(t *testing.T) {
	var termios unix.Termios
	newTermios := newTermios(termios)

	if &termios == &newTermios {
		t.Error("termios should be copied")
	}
}
