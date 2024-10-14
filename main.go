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
	fmt.Fprintf(os.Stderr, "%s^ ", strings.Repeat(" ", loc))
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

//
// Tokenizer
//

type (
	TokenKind int

	Token struct {
		kind  TokenKind
		index int
		val   int
		str   string
		len   int
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

		if strings.HasPrefix(p[i:], "==") || strings.HasPrefix(p[i:], "!=") ||
			strings.HasPrefix(p[i:], "<=") || strings.HasPrefix(p[i:], ">=") {
			tokens = append(tokens, Token{
				kind: tkReserved, index: i,
				str: p[i : i+2], len: 2,
			})
			i++
			continue
		}

		if strings.ContainsAny(string(c), "+-*/()<>") {
			tokens = append(tokens, Token{
				kind: tkReserved, index: i,
				str: string(c), len: 1,
			})
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
			errorAt(i, "invalid token")
		}
		tokens = append(tokens, Token{
			kind: tkNum, index: index,
			str: s, len: len(s), val: n,
		})
		i--
	}

	tokens = append(tokens, Token{kind: tkEof, index: len(p)})
	token = &tokens[0]
}

func advance() {
	tokens = tokens[1:]
	token = &tokens[0]
}

func consume(op string) bool {
	if token.kind != tkReserved || token.str != op {
		return false
	}
	advance()
	return true
}

func expect(op string) {
	if !consume(op) {
		errorAt(token.index, "expected %s", op)
	}
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

//
// Code generator
//

type (
	NodeKind int

	Node struct {
		kind     NodeKind
		lhs, rhs *Node
		val      int
	}
)

const (
	ndAdd NodeKind = iota
	ndSub
	ndMul
	ndDiv
	ndEq
	ndNe
	ndLt
	ndLe
	ndNum
)

func newNode(kind NodeKind, lhs, rhs *Node) *Node {
	return &Node{kind: kind, lhs: lhs, rhs: rhs}
}

func newNodeNum(val int) *Node {
	return &Node{kind: ndNum, val: val}
}

func primary() *Node {
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	return newNodeNum(expectNumber())
}

func unary() *Node {
	if consume("-") {
		return newNode(ndSub, newNodeNum(0), unary())
	} else if consume("+") {
		return unary()
	}

	return primary()
}

func mul() *Node {
	node := unary()

	for {
		if consume("*") {
			node = newNode(ndMul, node, unary())
		} else if consume("/") {
			node = newNode(ndDiv, node, unary())
		} else {
			return node
		}
	}
}

func add() *Node {
	node := mul()

	for {
		if consume("+") {
			node = newNode(ndAdd, node, mul())
		} else if consume("-") {
			node = newNode(ndSub, node, mul())
		} else {
			return node
		}
	}
}

func relational() *Node {
	node := add()

	for {
		if consume("<") {
			node = newNode(ndLt, node, add())
		} else if consume("<=") {
			node = newNode(ndLe, node, add())

		// opposite equivalents
		} else if consume(">") {
			node = newNode(ndLt, add(), node)
		} else if consume(">=") {
			node = newNode(ndLe, add(), node)
		} else {
			return node
		}
	}
}

func equality() *Node {
	node := relational()

	for {
		if consume("==") {
			node = newNode(ndEq, node, relational())
		} else if consume("!=") {
			node = newNode(ndNe, node, relational())
		} else {
			return node
		}
	}
}

func expr() *Node { return equality() }

func gen(node *Node) {
	if node.kind == ndNum {
		fmt.Printf("\tpush %d\n", node.val)
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	fmt.Println("\tpop rdi")
	fmt.Println("\tpop rax")

	switch node.kind {
	case ndAdd:
		fmt.Println("\tadd rax, rdi")
	case ndSub:
		fmt.Println("\tsub rax, rdi")
	case ndMul:
		fmt.Println("\timul rax, rdi")
	case ndDiv:
		fmt.Println("\tcqo")
		fmt.Println("\tidiv rdi")
	case ndEq:
		fmt.Println("\tcmp rax, rdi")
		fmt.Println("\tsete al")
		fmt.Println("\tmovzb rax, al")
	case ndNe:
		fmt.Println("\tcmp rax, rdi")
		fmt.Println("\tsetne al")
		fmt.Println("\tmovzb rax, al")
	case ndLt:
		fmt.Println("\tcmp rax, rdi")
		fmt.Println("\tsetl al")
		fmt.Println("\tmovzb rax, al")
	case ndLe:
		fmt.Println("\tcmp rax, rdi")
		fmt.Println("\tsetle al")
		fmt.Println("\tmovzb rax, al")
	}

	fmt.Println("\tpush rax")
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
