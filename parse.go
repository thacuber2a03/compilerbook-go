package main

// Copyright (c) 2024 @thacuber2a03
// This software is released under the terms of the MIT License. See LICENSE for details.

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

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
	tkIdent
	tkNum
	tkEof
)

var (
	userInput string
	tokens    = []Token{}
	token     *Token
)

func errorAt(loc int, format string, a ...any) {
	fmt.Fprintln(os.Stderr, userInput)
	fmt.Fprintf(os.Stderr, "%s^ ", strings.Repeat(" ", loc))
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

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

		if strings.ContainsAny(string(c), "+-*/()<>=;") {
			tokens = append(tokens, Token{
				kind: tkReserved, index: i,
				str: string(c), len: 1,
			})
			continue
		}

		if 'a' <= c && c <= 'z' {
			tokens = append(tokens, Token{
				kind: tkIdent, index: i,
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

func consumeIdent() (*Token, bool) {
	if token.kind != tkIdent {
		return nil, false
	}
	t := token
	advance()
	return t, true
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

type (
	NodeKind int

	Node struct {
		kind     NodeKind
		lhs, rhs *Node
		val      int // only for ndNum
		offset   int // only for ndLvar
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

	ndAssign
	ndLvar
	ndNum
)

var (
	code = make([]*Node, 0, 100)
)

func newNode(kind NodeKind, lhs, rhs *Node) *Node {
	return &Node{kind: kind, lhs: lhs, rhs: rhs}
}

func newNodeNum(val int) *Node { return &Node{kind: ndNum, val: val} }

// primary = num | ident | "(" expr ")"
func primary() *Node {
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	if tok, ok := consumeIdent(); ok {
		return &Node{
			kind:   ndLvar,
			offset: int(tok.str[0]-'a'+1) * 8,
		}
	}

	return newNodeNum(expectNumber())
}

// unary = ("+" | "-")? primary
func unary() *Node {
	if consume("-") {
		return newNode(ndSub, newNodeNum(0), unary())
	} else if consume("+") {
		return unary()
	}

	return primary()
}

// mul = unary ("*" unary | "/" unary)*
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

// add = mul ("+" mul | "-" mul)*
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

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *Node {
	node := add()

	for {
		if consume("<") {
			node = newNode(ndLt, node, add())
		} else if consume("<=") {
			node = newNode(ndLe, node, add())
		} else if consume(">") {
			node = newNode(ndLt, add(), node)
		} else if consume(">=") {
			node = newNode(ndLe, add(), node)
		} else {
			return node
		}
	}
}

// equality = relational ("==" relational | "!=" relational)*
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

// assign = equality ("=" assign)?
func assign() *Node {
	node := equality()
	if consume("=") {
		node = newNode(ndAssign, node, assign())
	}
	return node
}

// expr = assign
func expr() *Node { return assign() }

// stmt = expr ";"
func stmt() *Node {
	node := expr()
	expect(";")
	return node
}

// program = stmt*
func program() {
	for !atEof() {
		code = append(code, stmt())
	}
}
