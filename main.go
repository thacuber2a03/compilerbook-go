package main

import (
	"fmt"
	"os"
	"regexp"
)

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		die("usage: 9cc [input]")
	}

	p := os.Args[1]

	ws := regexp.MustCompile(`\b`).Split(p, -1)

	fmt.Println(`.intel_syntax noprefix`)
	fmt.Println(`.globl main`)
	fmt.Println(`main:`)
	fmt.Println("\tmov rax, ", ws[0])

	for i := 1; i < len(ws); i++ {
		w := ws[i]
		switch w[0] {
		case '+':
			i++
			fmt.Println("\tadd rax, ", ws[i])
		case '-':
			i++
			fmt.Println("\tsub rax, ", ws[i])
		default:
			die(`unexpected character '%c'`, w[0])
		}
	}

	fmt.Println("\tret")
}
