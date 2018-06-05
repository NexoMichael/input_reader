// +build darwin dragonfly freebsd netbsd openbsd

package input_reader

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA
const ioctlWriteTermios = unix.TIOCSETA
