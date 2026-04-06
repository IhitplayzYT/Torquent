package main

import (
	"fmt"
	"os"
)

type T_config struct {
	announce      string
	announce_list [][]string
	name          string
	piece_len     int
	length        int64
	pieces        [][]byte
	files         []Fmeta
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
	k, ok := announce.(string)
	if !ok {
		fmt.Println("Corrupt map announce cant be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.announce = k

	announce_list, ok := cfg["announce-list"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Announce List\"")
		os.Exit(int(E_FILE))
	}
	a, ok := announce_list.([][]string)
	if !ok {
		fmt.Println("Corrupt map announce_list cant be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.announce_list = a

	info, ok := cfg["info"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Info\"")
		os.Exit(int(E_FILE))
	}

	// TODO:  Get the info struct details
	// FIXME:

	return ret
}
