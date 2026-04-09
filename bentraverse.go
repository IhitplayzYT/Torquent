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
		os.Exit(int(E_NOTFOUND))
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

type Peer_v4 struct {
	ip   [4]byte
	port [2]byte
}

type Peer_v6 struct {
	ip [16]byte
}

type Peers struct {
	peers_v4 []Peer_v4
	peers_v6 []Peer_v6
	interval int
}

func parse_resp_dict(dict map[string]any) Peers {
	if fail, err := dict["failure reason"]; !err {
		fmt.Printf("ERROR IN RESPONSE MAP: %v\n", string(fail.([]byte)))
		os.Exit(int(E_HTTP))
	}

	interv, err := dict["interval"]
	if err {
		fmt.Println("No field interval found in Map")
		os.Exit(int(E_NOTFOUND))
	}
	iv, err := interv.(int)
	if err {
		fmt.Println("Unable to typecast interval to int")
		os.Exit(int(E_BADRES))
	}
	ret := Peers{
		interval: iv,
	}

	peers_v4, err := dict["peers"]
	if err {
		fmt.Println("No field peers found in Map")
		os.Exit(int(E_NOTFOUND))
	}
	pv4, err := peers_v4.([]byte)
	if err {
		fmt.Println("Unable to typecast peers to bytes")
		os.Exit(int(E_BADRES))
	}

	if len(pv4)%6 != 0 {
		fmt.Println("Corrupt bytes for peers")
		os.Exit(int(E_FILE))

	}
	lv4 := len(pv4) % 6
	v4 := make([]Peer_v4, 0)
	for i := 0; i < lv4; i++ {
		v4 = append(v4, Peer_v4{
			ip:   [4]byte(pv4[6*i : 6*i+4]),
			port: [2]byte(pv4[6*i+4 : 6*i+6]),
		})
	}
	ret.peers_v4 = v4

	peers_v6, err := dict["peers6"]
	if err {
		return ret
	}
	pv6, err := peers_v6.([]byte)
	if err {
		fmt.Println("Unable to typecast peers6 to bytes")
		os.Exit(int(E_BADRES))
	}

	if len(pv6)%18 != 0 {
		ret.peers_v4 = v4
		fmt.Println("Corrupt bytes for peers")
		os.Exit(int(E_FILE))

	}
	lv6 := len(pv6)
	v6 := make([]Peer_v6, 0)
	for i := 0; i < lv6; i++ {
		v6 = append(v6, Peer_v6{
			ip: [16]byte(pv6[18*i : 18*i+18]),
		})
	}

	ret.peers_v6 = v6

	return ret

}
