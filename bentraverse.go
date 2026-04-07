package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"os"
)

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
		indent(ident)
		fmt.Println("}")
	case BINT:
		indent(ident)
		fmt.Println("int:", node.Int)
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

func findinfo(node *BNode) *BNode {
	if node.Type != BDICT {
		fmt.Println("Torrent traverse failed due to corrupt torrent root")
		os.Exit(int(E_FILE))
	}
	nd, ok := node.Dict["info"]
	if !ok {
		fmt.Println("Torrent traverse failed due to missing info root")
		os.Exit(int(E_FILE))
	}
	return nd
}

func (self *Torrent) Merkle_root(root *BNode) (v1_hash [20]byte, v2_hash [32]byte, is_v1 bool, is_v2 bool) {
	info := findinfo(root)
	data := self.doc[info.span.strt:info.span.end]
	mv, hasMV := info.Dict["meta version"]
	isV2 := hasMV && mv.Int == 2
	_, hasPieces := info.Dict["pieces"]
	_, hasLen := info.Dict["length"]
	_, hasfiles := info.Dict["files"]
	isV1 := !isV2 || hasPieces || hasLen || hasfiles
	if isV1 {
		v1_hash = sha1.Sum(data)
	}
	if isV2 {
		v2_hash = sha256.Sum256(data)
	}
	return v1_hash, v2_hash, isV1, isV2
}
