package main

import (
	"bytes"
	"fmt"
	"os"
)

type Btype int

const (
	BINT Btype = iota
	BSTR
	BLIST
	BDICT
)

type Span struct {
	strt, end int
}

type BNode struct {
	Type Btype
	Int  int
	Str  []byte
	List []*BNode
	Dict map[string]*BNode
	span Span
}

type Torrent struct {
	fname string
	doc   []byte
	cur   int
	node  *BNode
}

func (t *Torrent) is_open() bool {
	return len(t.doc) != 0
}

func (t *Torrent) open_file() {
	f, err := os.ReadFile(t.fname)
	if err != nil {
		fmt.Println("Unable to open the file")
		os.Exit(int(E_IO))
	}
	t.doc = f
}

func init_torrent(name string) Torrent {
	self := Torrent{fname: name}
	self.open_file()

	if len(self.doc) <= 1 {
		fmt.Println("Truncated/Corrupt torrent file provided")
		os.Exit(int(E_FILE))
	}
	self.cur = 0
	self.node = nil
	return self
}

func (self *Torrent) Parse() *BNode {
	if self.cur >= len(self.doc) {
		fmt.Println("EOF")
		os.Exit(int(E_FILE))
	}
	strt := self.cur
	switch self.doc[self.cur] {
	case 'd':
		self.cur++
		return &BNode{Type: BDICT, Dict: self.eval_dict(), span: Span{strt, self.cur}}
	case 'l':
		self.cur++
		return &BNode{Type: BLIST, List: self.eval_list(), span: Span{strt, self.cur}}
	case 'i':
		self.cur++
		return &BNode{Type: BINT, Int: self.eval_int(), span: Span{strt, self.cur}}
	default:
		return &BNode{Type: BSTR, Str: self.eval_bstr(), span: Span{strt, self.cur}}
	}
}

// d....e
func (self *Torrent) eval_dict() map[string]*BNode {
	ret := make(map[string]*BNode)
	var prev []byte = nil

	for self.cur < len(self.doc) && self.doc[self.cur] != 'e' {
		kb := self.eval_bstr()
		if prev != nil && bytes.Compare(kb, prev) <= 0 {
			fmt.Println("Parsing failed @ Dict due to unsorted Key-Value")
			os.Exit(int(E_FILE))
		}
		prev = kb
		ret[string(kb)] = self.Parse()
	}

	if self.cur >= len(self.doc) || self.doc[self.cur] != 'e' {
		fmt.Println("Parsing failed @ Dict due to Corrupt/Invalid file")
		os.Exit(int(E_FILE))
	}
	self.cur++
	return ret
}

// l....e
func (self *Torrent) eval_list() []*BNode {
	k := make([]*BNode, 0)

	for self.cur < len(self.doc) && self.doc[self.cur] != 'e' {
		k = append(k, self.Parse())
	}

	if self.cur >= len(self.doc) || self.doc[self.cur] != 'e' {
		fmt.Println("Parsing failed @ List due to Corrupt/Invalid file")
		os.Exit(int(E_FILE))
	}
	self.cur++
	return k
}

// 5:hello
func (self *Torrent) eval_bstr() []byte {
	strt := self.cur
	for self.cur < len(self.doc) && self.doc[self.cur] != ':' {
		self.cur++
	}
	if self.cur >= len(self.doc) {
		fmt.Println("Parsing failed @ Str due to Corrupt/Invalid file")
		os.Exit(int(E_FILE))
	}

	le := self.doc[strt:self.cur]

	if len(le) > 1 && le[0] == '0' {
		fmt.Println("Parsing failed @ Str due to leading zero in length")
		os.Exit(int(E_FILE))
	}

	l := 0
	for _, v := range le {
		if v < '0' || v > '9' {
			fmt.Println("Parsing failed @ Str due to invalid string length")
			os.Exit(int(E_FILE))
		}
		l = l*10 + int(v-'0')
	}
	self.cur++

	if self.cur+l > len(self.doc) {
		fmt.Println("Parsing failed @ Str due to Corrupt/Invalid file")
		os.Exit(int(E_FILE))
	}

	ret := self.doc[self.cur : self.cur+l]
	self.cur += l
	return ret
}

func (self *Torrent) eval_int() int {

	neg := false
	if self.cur < len(self.doc) && self.doc[self.cur] == '-' {
		neg = true
		self.cur++
		if self.cur < len(self.doc) && self.doc[self.cur] == '0' {
			fmt.Println("Parsing failed @ Int due to invalid: -0")
			os.Exit(int(E_FILE))
		}
	}

	if self.cur >= len(self.doc) || self.doc[self.cur] < '0' || self.doc[self.cur] > '9' {
		fmt.Println("Parsing failed @ Int due to invalid Int")
		os.Exit(int(E_FILE))
	}

	if self.doc[self.cur] == '0' && self.cur+1 < len(self.doc) && self.doc[self.cur+1] != 'e' {
		fmt.Println("Parsing failed @ Int due to leading zeroes")
		os.Exit(int(E_FILE))
	}

	val := 0
	for self.cur < len(self.doc) && self.doc[self.cur] >= '0' && self.doc[self.cur] <= '9' {
		val = val*10 + int(self.doc[self.cur]-'0')
		self.cur++
	}

	if self.cur >= len(self.doc) || self.doc[self.cur] != 'e' {
		fmt.Println("Parsing failed @ Int due to Corrupt/Invalid file")
		os.Exit(int(E_FILE))
	}
	self.cur++ // consume 'e'

	if neg {
		return -val
	}
	return val
}
