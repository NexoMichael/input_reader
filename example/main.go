package main

import (
	"fmt"
	"os"
	"syscall"

	input "github.com/NexoMichael/inputreader"
)

func main() {
	reader, _ := input.NewInputReader(os.Stdin)
	defer func() {
		reader.Close()
	}()

	var b [1]byte
	for {
		n, err := reader.Read(b[:])
		if err != nil || n == 0 {
			return
		}

		fmt.Printf(" - code is %d\n\r", b[0])

		switch syscall.Signal(b[0]) {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		}
	}
}
