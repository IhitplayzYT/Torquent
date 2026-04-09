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
		query := ""
		if is_v1 {
			query = cfg.mk_v1query(v1)
		}
		if is_v2 {
			query = cfg.mk_v2query(v2)
		}
		fmt.Println(query)
		WRN(query)
	}

}
