package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var userInput string

func errorAt(loc int, format string, a ...any) {
	fmt.Fprintln(os.Stderr, userInput)
	fmt.Fprintf(os.Stderr, "%s^ ",  strings.Repeat(" ", loc))
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

type (
	TokenKind int

	Token struct {
		kind  TokenKind
		val   int
		str   string
		index int
	}
)

const (
	tkReserved TokenKind = iota
	tkNum
	tkEof
)

var (
	tokens = []Token{}
	token  *Token
)

func tokenize() {
	p := userInput

	for i := 0; i < len(p); i++ {
		c := p[i]
		if unicode.IsSpace(rune(c)) {
			continue
		}

		if c == '+' || c == '-' {
			tokens = append(tokens, Token{kind: tkReserved, str: string(c), index: i})
			continue
		}

		index := i
		var num strings.Builder
		for i < len(p) && unicode.IsDigit(rune(p[i])) {
			num.WriteByte(p[i])
			i++
		}
		s := num.String()

		n, e := strconv.Atoi(s)
		if e != nil {
			errorAt(i, "cannot tokenize")
		}
		tokens = append(tokens, Token{kind: tkNum, str: s, val: n, index: index})
		i--
	}

	tokens = append(tokens, Token{kind: tkEof, index: len(p)})
	token = &tokens[0]
}

func advance() {
	tokens = tokens[1:]
	token = &tokens[0]
}

func consume(op byte) bool {
	if token.kind != tkReserved || token.str[0] != op {
		return false
	}
	advance()
	return true
}

func expect(op byte) {
	if token.kind != tkReserved || token.str[0] != op {
		errorAt(token.index, "expected %c", op)
	}
	advance()
}

func expectNumber() int {
	if token.kind != tkNum {
		errorAt(token.index, "expected a number")
	}

	val := token.val
	advance()
	return val
}

func atEof() bool { return token.kind == tkEof }

func main() {
	if len(os.Args) != 2 {
		die("usage: 9cc [input]")
	}

	userInput = os.Args[1]
	tokenize()

	fmt.Println(`.intel_syntax noprefix
.globl main
main:`)

	fmt.Printf("\tmov rax, %d\n", expectNumber())

	for !atEof() {
		if consume('+') {
			fmt.Printf("\tadd rax, %d\n", expectNumber())
			continue
		}

		expect('-')
		fmt.Printf("\tsub rax, %d\n", expectNumber())
	}

	fmt.Println("\tret")
}
