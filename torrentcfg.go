package main

import (
	"fmt"
	"os"
)

type T_config struct {
	announce      string     //
	announce_list [][]string //
	name          string     //
	piece_len     int        //
	length        int64      //
	pieces        [][20]byte //
	n_pieces      int        //
	files         []Fmeta    //
}

type Fmeta struct {
	path   []string
	length int64
	offset int64
}

func get_meta(node *BNode, cfg map[string]any) T_config {
	ret := T_config{}
	announce, ok := cfg["announce"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Announce\"")
		os.Exit(int(E_FILE))
	}
	a, ok := announce.(string)
	if !ok {
		fmt.Println("Corrupt map announce cant be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.announce = a

	announce_list, ok := cfg["announce-list"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Announce List\"")
		os.Exit(int(E_FILE))
	}
	b, ok := announce_list.([][]string)
	if !ok {
		fmt.Println("Corrupt map announce_list cant be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.announce_list = b

	info, ok := cfg["info"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Info\"")
		os.Exit(int(E_FILE))
	}

	info_map, ok := info.(map[string]any)
	if !ok {
		fmt.Println("Corrupt map info can't be typecasted")
		os.Exit(int(E_FILE))
	}

	c, ok := info_map["length"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"lenght\"")
		os.Exit(int(E_FILE))
	}

	d, ok := c.(int64)
	if !ok {
		fmt.Println("Corrupt info_map length can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.length = d

	e, ok := info_map["name"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"name\"")
		os.Exit(int(E_FILE))
	}

	f, ok := e.(string)
	if !ok {
		fmt.Println("Corrupt info_map name can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.name = f

	g, ok := info_map["piece length"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"piece length\"")
		os.Exit(int(E_FILE))
	}

	h, ok := g.(int)
	if !ok {
		fmt.Println("Corrupt info_map length can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.piece_len = h

	i, ok := info_map["files"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"files\"")
		os.Exit(int(E_FILE))
	}

	j, ok := i.([]Fmeta)
	if !ok {
		fmt.Println("Corrupt info_map files can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.files = j

	l, ok := info_map["pieces"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"pieces\"")
		os.Exit(int(E_FILE))
	}

	m, ok := l.([][20]byte)
	if !ok {
		fmt.Println("Corrupt info_map pieces can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.pieces = m
	ret.n_pieces = len(ret.pieces)

	return ret
}
