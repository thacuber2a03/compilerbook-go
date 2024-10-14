package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func error(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

type (
	TokenKind int
	Token struct {
		kind TokenKind
		val int
		str string
	}
)

const (
	tkReserved TokenKind = iota
	tkNum
	tkEof
)

var (
	tokens = []Token{}
	token *Token
)

func tokenize(p string) {
	for i := 0; i < len(p); i++ {
		c := p[i]
		if unicode.IsSpace(rune(c)) {
			continue
		}

		if c == '+' || c == '-' {
			tokens = append(tokens, Token{kind: tkReserved, str: string(c)})
			continue
		}

		var num strings.Builder
		for i < len(p) && unicode.IsDigit(rune(p[i])) {
			num.WriteByte(p[i])
			i++
		}
		i--
		s := num.String()

		n, e := strconv.Atoi(s)
		if e != nil {
			error("cannot tokenize")
		}
		tokens = append(tokens, Token{kind: tkNum, str: s, val: n})
	}

	tokens = append(tokens, Token{kind: tkEof})
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
		error("expected %c", op)
	}
	advance()
}

func expectNumber() int {
	if token.kind != tkNum {
		error("expected a number")
	}

	val := token.val
	advance()
	return val
}

func atEof() bool { return token.kind == tkEof }

func main() {
	if len(os.Args) != 2 {
		error("usage: 9cc [input]")
	}

	tokenize(os.Args[1])

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
