package main

import (
	"fmt"
	"os"
)

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		die("usage: 9cc [input]")
	}

	userInput = os.Args[1]
	tokenize()
	node := expr()

	fmt.Println(`.intel_syntax noprefix
.globl main
main:`)

	gen(node)

	fmt.Println("\tpop rax")
	fmt.Println("\tret")
}
