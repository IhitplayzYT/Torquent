package main

import (
	"fmt"
	"math/rand"
)

var PEER_IDS map[string]bool = make(map[string]bool, 0)

const CLIENT_ID string = "-RU0001-"
const RANDSTR string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Progress struct {
	uploaded, dowloaded, left int
}

var progress Progress = Progress{uploaded: 0, dowloaded: 0, left: 0}

func gen_pid() string {
	l := len(RANDSTR)
	UNIQ := true
	for UNIQ {
		i := 0
		ret := ""
		for i < 12 {
			ret += string(RANDSTR[rand.Int()%l])
			i += 1
		}
		if _, ok := PEER_IDS[ret]; !ok {
			PEER_IDS[ret] = true
			return ret
		}
	}
	return "XXXXXXXXXXXX"
}

func (t T_config) fmt_hash(hash [20]byte) string {
	ret := "%"
	for _, v := range hash {
		ret += fmt.Sprintf("%x%%", v)
	}
	return ret[:len(ret)-1]
}

// GET /announce?info_hash=<urlencoded>&peer_id=<id>&port=6881&uploaded=0&downloaded=0&left=1048576&compact=1

func (t T_config) mk_v1query(hash [20]byte) string {
	return fmt.Sprintf(" GET /announce?info_hash=%v&peer_id=%v&port=6881&uploaded=%v&downloaded=%v&left=%v&compact=1", t.fmt_hash(hash), CLIENT_ID+gen_pid(), progress.uploaded, progress.dowloaded, progress.left)
}

func (t T_config) mk_v2query(hash [32]byte) string {
	return fmt.Sprintf("%v", hash)
}
