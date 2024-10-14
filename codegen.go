package main

// this file is pretty small...

import "fmt"

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
