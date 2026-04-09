package main

import (
	"fmt"
	"os"
)

type T_config struct {
	announce      string     //
	announce_list [][]string //
	name          string     //
	piece_len     int64      //
	length        int64      //
	pieces        [][]byte   //
	n_pieces      int        //
	files         []Fmeta    //
}

func (t T_config) show() {
	fmt.Printf("T_CONFIG:\nannounce: %v\nannounce-list: %v\npieces[%v pieces of %v length]: %v\nFileName: (%v bytes)%v\nFiles: %v\n", t.announce, t.announce_list, t.n_pieces, t.piece_len, t.pieces, t.length, t.name, t.fmt_files())
}

func (t T_config) fmt_files() string {
	ret := "[\n"
	for _, v := range t.files {
		ret += fmt.Sprintf("%v of size %v bytes, offset at %v bytes\n", v.path, v.length, v.offset)
	}
	ret += "]\n"
	return ret
}

type Fmeta struct {
	path   []string
	length int64
	offset int64
}

func get_meta(cfg map[string]any) T_config {
	ret := T_config{}
	announce, ok := cfg["announce"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Announce\"")
		os.Exit(int(E_FILE))
	}
	a, ok := announce.([]byte)
	if !ok {
		fmt.Println("Corrupt map announce cant be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.announce = string(a)
	announce_list, ok := cfg["announce-list"]
	if !ok {
		fmt.Println("Corrupt torrent unable to find \"Announce List\"")
		os.Exit(int(E_FILE))
	}
	b, ok := announce_list.([]any)
	if !ok {
		fmt.Println("Corrupt map announce_list cant be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.announce_list = make([][]string, len(b))

	for i, tier := range b {
		tierList, ok := tier.([]any)
		if !ok {
			fmt.Println("invalid tier")
			os.Exit(1)
		}

		ret.announce_list[i] = make([]string, len(tierList))

		for j, tracker := range tierList {
			b, ok := tracker.([]byte)
			if !ok {
				fmt.Println("tracker not []byte")
				os.Exit(1)
			}

			ret.announce_list[i][j] = string(b)
		}
	}

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

	d, ok := c.(int)
	if !ok {
		fmt.Println("Corrupt info_map length can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.length = int64(d)

	e, ok := info_map["name"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"name\"")
		os.Exit(int(E_FILE))
	}
	f, ok := e.([]byte)
	if !ok {
		fmt.Println("Corrupt info_map name can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.name = string(f)

	g, ok := info_map["piece length"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"piece length\"")
		os.Exit(int(E_FILE))
	}

	h, ok := g.(int)
	if !ok {
		fmt.Println("Corrupt info_map piece length can't be typecasted")
		os.Exit(int(E_FILE))
	}
	ret.piece_len = int64(h)

	l, ok := info_map["pieces"]
	if !ok {
		fmt.Println("Corrupt info_map unable to find field \"pieces\"")
		os.Exit(int(E_FILE))
	}

	m, ok := l.([]byte)
	if !ok {
		fmt.Println("Corrupt info_map pieces can't be typecasted")
		os.Exit(int(E_FILE))
	}
	if len(m)%20 != 0 {
		fmt.Println("Corrupt info_map each piece has to be 20 bytes")
		os.Exit(int(E_FILE))
	}
	ret.n_pieces = int(len(m) / 20)
	ret.pieces = make([][]byte, ret.n_pieces)
	for idx := 0; idx < ret.n_pieces; idx++ {
		ret.pieces[idx] = m[idx*20 : idx*20+20]
	}

	//  SPECIAL CARE
	i, ok := info_map["files"]
	if ok {
		files_raw, ok := i.([]any)
		if !ok {
			fmt.Println("Corrupt info_map files not []any")
			os.Exit(int(E_FILE))
		}

		var offset int64 = 0
		ret.files = make([]Fmeta, len(files_raw))

		for idx, file := range files_raw {
			file_map, ok := file.(map[string]any)
			if !ok {
				fmt.Println("Corrupt file entry not map")
				os.Exit(int(E_FILE))
			}

			l_raw, ok := file_map["length"]
			if !ok {
				fmt.Println("Missing file length")
				os.Exit(int(E_FILE))
			}

			length, ok := l_raw.(int64)
			if !ok {
				fmt.Println("Invalid file length type")
				os.Exit(int(E_FILE))
			}

			p_raw, ok := file_map["path"]
			if !ok {
				fmt.Println("Missing file path")
				os.Exit(int(E_FILE))
			}

			p_list, ok := p_raw.([]any)
			if !ok {
				fmt.Println("Invalid path type")
				os.Exit(int(E_FILE))
			}

			path := make([]string, len(p_list))
			for i, seg := range p_list {
				s, ok := seg.(string)
				if !ok {
					fmt.Println("Invalid path segment")
					os.Exit(int(E_FILE))
				}
				path[i] = s
			}

			ret.files[idx] = Fmeta{
				path:   path,
				length: length,
				offset: offset,
			}

			offset += length
		}

	} else {
		ret.files = []Fmeta{
			{
				path:   []string{ret.name},
				length: ret.length,
				offset: 0,
			},
		}
	}

	progress.left = int(ret.length)
	return ret
}
