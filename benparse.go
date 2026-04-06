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

func (t *Torrent) open_file() {
	f, err := os.ReadFile(t.fname)

	if err != nil {
		fmt.Println("Unable to open the file")
		os.Exit(int(E_IO))
	}
	if len(t.doc) == 0 {
		t.doc = f
	}
}

func (t *Torrent) is_open() bool {
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
		os.Exit(int(E_FILE))
	}
	self.cur = 0
	self.end = l - 1
	self.node = nil
	return self
}

func (self *Torrent) Parse() *BNode {
	if self.cur >= len(self.doc) {
		fmt.Println("EOF")
		os.Exit(int(E_FILE))
	}
	switch self.doc[self.cur] {
	case 'd':
		self.cur++
		return &BNode{Type: BDICT, Dict: self.eval_dict()}
	case 'l':
		self.cur++
		return &BNode{Type: BLIST, List: self.eval_list()}
	case 'i':
		return &BNode{Type: BINT, Int: self.eval_int()}
	default:
		return &BNode{Type: BSTR, Str: self.eval_bstr()}
	}
}

// d....e
func (self *Torrent) eval_dict() map[string]*BNode {

	ret := make(map[string]*BNode, 0)
	for self.cur < len(self.doc) && self.doc[self.cur] != 'e' {
		k, v := string(self.eval_bstr()), self.Parse()
		ret[k] = v
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

	if len(le) == 0 {
		fmt.Println("Parsing failed @ Str due to Empty File")
		os.Exit(int(E_FILE))
	}
	if len(le) > 1 && le[0] == '0' {
		fmt.Println("Parsing failed @ Str due to null-terminated beginning")
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

	if self.cur > len(self.doc) || l < 0 {
		fmt.Println("Parsing failed @ Str due to Corrupt/Invaldi file")
		os.Exit(int(E_FILE))
	}

	if self.cur+l > len(self.doc) {
		fmt.Println("Parsing failed @ Str due to Corrupt/Invalid file")
		os.Exit(int(E_FILE))
	}

	ret := self.doc[self.cur : self.cur+l]
	self.cur += l
	return ret
}

// i1893179e
func (self *Torrent) eval_int() int {
	if self.doc[self.cur] == 'i' {
		self.cur++
	}

	neg := false
	if self.doc[self.cur] == '-' {
		neg = true
		self.cur++
		if self.cur < len(self.doc) && self.doc[self.cur] == '0' {
			fmt.Println("Parsing failed @ Int due to invalid length: -0")
			os.Exit(int(E_FILE))
		}
	}

	if self.cur >= len(self.doc) || self.doc[self.cur] < '0' || self.doc[self.cur] > '9' {
		fmt.Println("Parsing failed @ Int due to invalid Int")
		os.Exit(int(E_FILE))
	}

	if self.doc[self.cur] == '0' && self.cur+1 < len(self.doc) && self.doc[self.cur+1] != 'e' {
		fmt.Println("Parsing failed @ Int due to Int with leading zeroes")
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
	self.cur++

	if neg {
		return -val
	}
	return val
}
