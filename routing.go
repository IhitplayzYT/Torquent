package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var PEER_IDS map[string]bool = make(map[string]bool, 0)

const CLIENT_ID string = "-RU0001-"
const RANDSTR string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const PORT int = 6881

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

func (t T_config) mk_v1query(hash [20]byte) (string, string) {
	pid := CLIENT_ID + gen_pid()
	return fmt.Sprintf("%v?info_hash=%v&peer_id=%v&port=%v&uploaded=%v&downloaded=%v&left=%v&compact=1", t.announce, t.fmt_hash(hash), url.QueryEscape(pid), PORT, progress.uploaded, progress.dowloaded, progress.left), pid
}

func (t T_config) mk_v2query(hash [32]byte) string {
	return fmt.Sprintf("%v", hash)
}

func announce(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Duration(time.Duration.Milliseconds(10000)),
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ru-torrent/0.1")
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request Failed with status code: %v\n", resp.StatusCode)
	}
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
