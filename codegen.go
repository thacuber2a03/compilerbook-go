package main

// Copyright (c) 2024 @thacuber2a03
// This software is released under the terms of the MIT License. See LICENSE for details.

// this file is pretty small...

import "fmt"

func genLval(node *Node) {
	if node.kind != ndLvar {
		die("left side of assignment isn't a variable")
	}

	fmt.Println("\tmov rax, rbp")
	fmt.Printf("\tsub rax, %d\n", node.offset)
	fmt.Println("\tpush rax")
	fmt.Println()
}

func gen(node *Node) {
	switch node.kind {
	case ndNum:
		fmt.Printf("\tpush %d\n\n", node.val)
		return
	case ndLvar:
		genLval(node)
		fmt.Println("\tpop rax")
		fmt.Println("\tmov rax, [rax]")
		fmt.Println("\tpush rax")
		fmt.Println()
		return
	case ndAssign:
		genLval(node.lhs)
		gen(node.rhs)

		fmt.Println("\tpop rdi")
		fmt.Println("\tpop rax")
		fmt.Println("\tmov [rax], rdi")
		fmt.Println("\tpush rdi")
		fmt.Println()
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
	fmt.Println()
}
