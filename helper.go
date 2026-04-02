package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	dbg   bool
	files []string
	multi bool
	odir  string
	idir  string
}

func DEBUG_STR() {
	fmt.Println("./Torquent [-d|-t|-o|-i|-h] <TORRENT_FILES> ... \n  -d: Enables debug mode\n  -t: Enables multi threaded mode\n  -i=/path/to/input/dir: Provide the input directory to search for torrent file\n  -o=/path/to/output/dir: Provide the output directory for downloaded file\n  -h: Displays the help menu")
}

func parse_cli(cwd string) CLI {
	cli := CLI{
		dbg:   false,
		files: make([]string, 0),
		multi: false,
		odir:  cwd,
		idir:  cwd,
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("No Home directory found for the user")
	}

	args := os.Args[1:]
	if len(args) < 2 {
		DEBUG_STR()
		os.Exit(0)
	}
	for _, v := range args {
		if strings.HasPrefix(v, "-o=") {
			v = v[3:]
			v = strings.ReplaceAll(v, "~", home)
			dir, err := os.Stat(v)
			if errors.Is(err, os.ErrNotExist) {
				err2 := os.Mkdir(v, 0755)
				if err2 != nil {
					fmt.Println("Can't create the directory with path: ", v)
					os.Exit(-1)
				}
				idx := strings.LastIndex(v, "/")
				fmt.Println("Created directory ", v[idx+1:], "in path: ", v[:idx])
			} else if !dir.IsDir() {
				fmt.Println("The output path is not a directory: ", v)
			}
			cli.odir = v
			continue
		}
		if strings.HasPrefix(v, "-i=") {
			v = v[3:]
			v = strings.ReplaceAll(v, "~", home)
			dir, err := os.Stat(v)
			if errors.Is(err, os.ErrNotExist) {
				err2 := os.Mkdir(v, 0755)
				if err2 != nil {
					fmt.Println("Can't create the directory with path: ", v)
					os.Exit(-1)
				}
				idx := strings.LastIndex(v, "/")
				fmt.Println("Created directory ", v[idx+1:], "in path: ", v[:idx])
			} else if !dir.IsDir() {
				fmt.Println("The input path is not a directory: ", v)
			}
			cli.idir = v
			continue
		}
		if strings.HasSuffix(v, ".torrent") {
			cli.files = append(cli.files, v)
			continue
		}

		switch v {
		case "-d":
			cli.dbg = true
		case "-t":
			cli.multi = true
		case "-h":
			DEBUG_STR()
			os.Exit(1)
		default:
			fmt.Println(v)
			DEBUG_STR()
			os.Exit(1)
		}

	}

	return cli
}

func If[T any](cond bool, tval T, fval T) T {
	if cond {
		return tval
	}
	return fval
}
