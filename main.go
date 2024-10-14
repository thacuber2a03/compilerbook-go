package main

// Copyright (c) 2024 @thacuber2a03
// This software is released under the terms of the MIT License. See LICENSE for details.

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
	program()

	fmt.Println(`.intel_syntax noprefix
.globl main
main:`)

	// allocate space for 26 variables (jeez)
	fmt.Println("\tpush rbp")
	fmt.Println("\tmov rbp, rsp")
	fmt.Println("\tsub rsp, 208")
	fmt.Println()

	for _, e := range code {
		gen(e)
		// pop result of expression from stack to avoid overflow
		fmt.Println("\tpop rax")
	}

	// result of last expression is in rax; it will be the return value
	fmt.Println("\n\tmov rsp, rbp")
	fmt.Println("\tpop rbp")
	fmt.Println("\tret")
}
