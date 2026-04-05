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
	Str  []byte
	List []*BNode
	Dict map[string]*BNode
}

type Torrent struct {
	fname    string
	doc      []byte
	cur, end int
	node     *BNode
}

func (t Torrent) open_file() {
	f, err := os.ReadFile(t.fname)

	if err != nil {
		fmt.Println("Unable to open the file")
		os.Exit(-2)
	}
	t.doc = If(len(t.doc) == 0, f, t.doc)
}

func (t Torrent) is_open() bool {
	return len(t.doc) != 0
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
	self.cur = 0
	self.end = l - 1
	self.node = nil
	return self
}

func (self Torrent) Parse() *BNode {
	if self.cur == self.end {
		fmt.Println("Finished Parsing")
	}
	if self.doc[self.cur] == 'd' {
		if self.doc[self.end] == 'e' {
			return &BNode{Type: BDICT, Dict: self.eval_dict()}
		} else {
			fmt.Println("Corrupt File")
			os.Exit(-4)
		}

	} else if self.doc[self.cur] == 'i' {
		node := &BNode{
			Type: BINT,
		}
		if self.doc[self.end] == 'e' {
			node.Int = self.eval_int()
		} else {
			fmt.Println("Corrupt File")
			os.Exit(-4)
		}
		return node
	} else if self.doc[self.cur] == 'l' {
		if self.doc[self.end] == 'e' {
			return &BNode{Type: BLIST, List: self.eval_list()}
		} else {
			fmt.Println("Corrupt File")
			os.Exit(-4)
		}

	} else if self.doc[self.cur] >= '0' && self.doc[self.cur] <= '9' {
		return &BNode{Type: BSTR, Str: self.eval_bstr()}
	}

	fmt.Println("Invalid token encountered")
	return nil

}

// d....e
func (self Torrent) eval_dict() map[string]*BNode {
	ret := make(map[string]*BNode, 0)
	for self.doc[self.cur] != 'e' {
		k, v := string(self.eval_bstr()), self.Parse()
		ret[k] = v
	}
	if self.doc[self.cur] != 'e' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}

	return ret
}

// l....e
func (self Torrent) eval_list() []*BNode {
	if self.doc[self.cur] == 'l' {
		self.cur += 1
	}

	k := make([]*BNode, 0)
	for self.cur < self.end && self.doc[self.cur] != 'e' {
		k = append(k, self.Parse())
	}

	if self.doc[self.cur] != 'e' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}

	return k
}

//  2989281:wnfioeiworijj

func (self Torrent) eval_bstr() []byte {
	strt := self.cur
	for self.doc[self.cur] != ':' && self.cur != self.end {
		self.cur += 1
	}
	if self.cur == self.end {
		fmt.Println("Corrupt Torrent File")
		os.Exit(-4)
	}
	le := self.doc[strt:self.cur]
	l := 0
	for _, v := range le {
		l = l*10 + int(v-'0')
	}

	self.cur += 1

	if self.cur == self.end || l < 0 {
		fmt.Println("Corrupt Torrent File")
		os.Exit(-4)
	}

	if self.cur+l > self.end {
		fmt.Println("Corrupt Torrent File")
		os.Exit(-4)
	}

	ret := self.doc[self.cur : self.cur+l]
	self.cur += l
	return ret

}

// i1893179e
// ^       ^
func (self Torrent) eval_int() int {
	if self.doc[self.cur] == 'i' {
		self.cur += 1
	}
	neg := If(self.doc[self.cur] == '-', true, false)
	val := 0
	for self.cur < self.end && (self.doc[self.cur] >= '0' && self.doc[self.cur] <= '9') {
		val = val*10 + int(self.doc[self.cur]-'0')
		self.cur += 1
	}

	if self.doc[self.cur] != 'e' {
		fmt.Println("Corrupt torrent file")
		os.Exit(-4)
	}

	return If(neg, -val, val)
}
