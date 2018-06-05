package main

import (
	"fmt"
	"os"

	input "github.com/NexoMichael/inputreader"
)

func main() {
	reader, _ := input.NewInputReader(os.Stdin)
	defer reader.Close()

	var b [1]byte
	for {
		n, err := reader.Read(b[:])
		if err != nil || n == 0 {
			return
		}

		fmt.Printf(" - code is %d\n\r", b[0])
	}
}
