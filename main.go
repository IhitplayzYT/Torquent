package main

import (
	"fmt"
	"os"
)

func main() {
	home, err := os.Getwd()
	if err != nil {
		fmt.Println("CRITICAL] No Current directory found!!")
	}
	clargs := parse_cli(home)
	if clargs.dbg {
		fmt.Printf("\nCLI {\ndbg: %s\nmulti: %s\ninputdir: %s\noutputdir: %s\nfiles: %v\n}\n", If(clargs.dbg, "True", "False"), If(clargs.multi, "Enabled", "Disabled"), clargs.idir, clargs.odir, clargs.files)
	}

}
