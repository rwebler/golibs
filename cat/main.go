package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		file, err := os.Open(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		io.Copy(os.Stdout, file)
	} else {
		fmt.Println("No filename")
		os.Exit(1)
	}
}
