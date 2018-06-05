package inputreader // import "github.com/NexoMichael/inputreader"

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TCGETS
const ioctlWriteTermios = unix.TCSETS
