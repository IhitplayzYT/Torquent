package main

import (
	"fmt"
	"os"
)

func init_cli() CLI {
	home, err := os.Getwd()
	if err != nil {
		fmt.Println("CRITICAL] No Current directory found!!")
		os.Exit(int(E_IO))
	}
	clargs := parse_cli(home)
	if clargs.dbg {
		fmt.Printf("\nCLI {\ndbg: %s\nmulti: %s\ninputdir: %s\noutputdir: %s\nfiles: %v\n}\n\n", If(clargs.dbg, "True", "False"), If(clargs.multi, "Enabled", "Disabled"), clargs.idir, clargs.odir, clargs.files)
	}
	return clargs
}

func main() {
	clargs := init_cli()
	for i, v := range clargs.files {
		torrconfig := init_torrent(v)
		if !torrconfig.is_open() {
			fmt.Println("Fopen has failed for", i, If(i%10 == 1, "st", If(i%10 == 2, "nd", If(i%10 == 3, "rd", "th"))), "file, named: ", v)
			os.Exit(int(E_IO))
		}
		node := torrconfig.Parse()
		if node == nil {
			fmt.Println("Parse failed!!")
			os.Exit(int(E_BEN))
		}
		if clargs.dbg {
			print_tree(node, 0)
		}
		dict := Traverse(node).(map[string]any)
		cfg := get_meta(dict)
		if clargs.dbg {
			cfg.show()
		}
		v1, v2, is_v1, is_v2 := torrconfig.Merkle_root(node)
		if is_v1 {
			query, poid := cfg.mk_v1query(v1)
			response, err := announce(query)
			if err != nil {
				fmt.Printf("Error in sending request: %v\n", err)
				os.Exit(int(E_HTTP))
			}
			resp_torrconfig := Torrent{
				doc:   response,
				fname: "",
				cur:   0,
				node:  nil,
			}
			resp_node := resp_torrconfig.Parse()
			if resp_node == nil {
				fmt.Println("Parse failed!!")
				os.Exit(int(E_BEN))
			}
			if clargs.dbg {
				print_tree(resp_node, 0)
			}

			r_dict := Traverse(resp_node).(map[string]any)
			peers := parse_resp_dict(r_dict)
			pid := make([]byte, 20)
			copy(pid, []byte(poid[:20]))
			get_pieces(peers, v1, [20]byte(pid), cfg.pieces)

		}
		if is_v2 {
			query := cfg.mk_v2query(v2)
			WRN(query)
		}
	}

}
