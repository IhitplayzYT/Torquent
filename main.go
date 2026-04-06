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
		fmt.Printf("\nCLI {\ndbg: %s\nmulti: %s\ninputdir: %s\noutputdir: %s\nfiles: %v\n}\n", If(clargs.dbg, "True", "False"), If(clargs.multi, "Enabled", "Disabled"), clargs.idir, clargs.odir, clargs.files)
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

	}

}
