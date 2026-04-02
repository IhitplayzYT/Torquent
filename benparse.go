package main

import (
	"fmt"
	"os"
)

type Parser struct {
	cur, end int
}
type Btype int

const (
	BINT Btype = iota
	BSTR
	BLIST
	BDICT
)

type BNode struct {
	Type Btype
	Int  int
	Str  string
	List []*BNode
	Dict map[string]*BNode
}

type Torrent struct {
	fname  string
	doc    []byte
	parser Parser
	node   *BNode
}

func (t Torrent) open_file() {
	f, err := os.ReadFile(t.fname)

	if err != nil {
		fmt.Println("Unable to open the file")
		os.Exit(-2)
	}
	t.doc = f
}

func (t Torrent) is_open() bool {
	return len(t.doc) == 0
}

func (t Torrent) idx() (int, int) {
	return t.parser.cur, t.parser.end
}

func init_torrent(name string) Torrent {
	self := Torrent{fname: name}

	if !self.is_open() {
		self.open_file()
	}
	l := len(self.doc)
	if l <= 1 {
		fmt.Println("Truncated/Corrupt torrent file provided")
		os.Exit(-3)
	}
	self.parser.cur = 0
	self.parser.end = l - 1
	self.node = nil
	return self
}

func (self Torrent) Parse() *BNode {
	if self.parser.cur == self.parser.end {
		fmt.Println("Finished Parsing")
	}
	l, r := self.idx()
	if self.doc[l] == 'd' {
		if self.doc[r] == 'e' {
			return self.eval_dict()
		}

	} else if self.doc[l] == 'i' {
		node := &BNode{
			Type: BINT,
		}
		if self.doc[r] == 'e' {
			node.Int = self.eval_int()
		} else {
			fmt.Println("Corrupt File")
			os.Exit(-4)
		}
		return node
	} else if self.doc[l] == 'l' {
		if self.doc[r] == 'e' {
			return self.eval_list()
		}

	} else if self.doc[l] >= 0 && self.doc[l] <= 9 {
		return self.eval_bstr()
	}

	fmt.Println("Invalid token encountered")
	return nil

}

func (self Torrent) eval_dict() map[string]*BNode {
	ret := make(map[string]*BNode, 0)
	l, r := self.idx()
	for self.doc[l] != 'e' {
		k, v := string(self.eval_bstr()), self.Parse()
		ret[k] = v
	}
	if self.doc[l] != 'e' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}

	self.parser.cur = l
	self.parser.end = r
	return ret
}

func (self Torrent) eval_list() []*BNode {

	self.parser.cur += 1
	l, r := self.idx()

	k := make([]*BNode, 0)
	for l < r && self.doc[l] != 'e' {
		k = append(k, self.Parse())
	}

	if self.doc[l] != 'e' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}

	self.parser.cur = l
	self.parser.end = r
	return k
}
func (self Torrent) eval_bstr() []byte {
	l, r := self.idx()
	if self.doc[l] == 's' {
		l += 1
	}
	node := &BNode{
		Type: BLIST,
	}

	neg := If(self.doc[l] == '-', true, false)
	val := 0
	for l < r && (self.doc[l] >= '0' && self.doc[l] <= '9') {
		val = val*10 + int(self.doc[l]-'0')
		l += 1
	}
	k := If(neg, -val, val)

	if self.doc[l] != ':' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}
	l += 1
	if l+k > r {
		fmt.Println("Corrupt torrent file")
		os.Exit(-5)
	}

	node.Str = string(self.doc[l : l+k+1])
	self.parser.cur = l + k
	self.parser.end = r
	return node
}

func (self Torrent) eval_int() int {
	l, r := self.idx()
	if self.doc[l] == 'i' {
		l += 1
	}
	neg := If(self.doc[l] == '-', true, false)
	val := 0
	for l < r && (self.doc[l] >= '0' && self.doc[l] <= '9') {
		val = val*10 + int(self.doc[l]-'0')
		l += 1
	}

	if self.doc[l] != 'e' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}
	self.parser.cur = l
	self.parser.end = r
	return If(neg, -val, val)
}
