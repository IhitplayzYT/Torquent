package main

import "fmt"

func Traverse(node *BNode) any {
	switch node.Type {
	case BDICT:
		ret := make(map[string]any, len(node.Dict))
		for k, v := range node.Dict {
			ret[k] = Traverse(v)
		}
		return ret
	case BINT:
		return node.Int
	case BLIST:
		ret := make([]any, len(node.List))
		for i, v := range node.List {
			ret[i] = Traverse(v)
		}
		return ret
	case BSTR:
		return node.Str
	default:
		return nil
	}
}

func print_tree(node *BNode, ident int) {
	indent := func(d int) {
		for i := 0; i < d; i++ {
			fmt.Print(" ")
		}
	}

	switch node.Type {
	case BDICT:
		indent(ident)
		fmt.Println("Dict: {")
		for k, v := range node.Dict {
			indent(ident + 1)
			fmt.Println("key:", k)
			print_tree(v, ident+2)
		}
		fmt.Println("}")
	case BINT:
		indent(ident)
		fmt.Print("int:", node.Int)
	case BLIST:
		indent(ident)
		fmt.Println("List: [")
		for _, v := range node.List {
			print_tree(v, ident+1)
		}
		fmt.Println("]")
	case BSTR:
		indent(ident)
		fmt.Printf("str(%d): %s\n", len(node.Str), node.Str)
	default:
		return
	}
}
