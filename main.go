package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: 9cc [input]\n")
		os.Exit(1)
	}

	v, e := strconv.Atoi(os.Args[1])
	if e != nil {
		panic(e)
	}

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
	mov rax, %d
	ret
`, v)
}
