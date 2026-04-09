package main

type Err int

const (
	E_IO Err = iota
	E_FS
	E_HELP
	E_CLI
	E_BEN
	E_FILE
	E_HTTP
	E_NOTFOUND
	E_BADRES
	E_NOPEER
)
