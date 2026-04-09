package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func get_pieces(peers Peers, hash [20]byte, peerid [20]byte) {
	if len(peers.peers_v4) == 0 && len(peers.peers_v6) == 0 {
		fmt.Println("No ipv4 & ipv6 peers")
		os.Exit(int(E_NOPEER))
	}

	if len(peers.peers_v4) != 0 {
		for i, v := range peers.peers_v4 {
			host := fmt.Sprintf("%s:%d", v.ip, v.port)
			conn, err := net.DialTimeout("tcp", host, 5*time.Second)
			if err != nil {
				fmt.Println("Error in Setting up TCP with ", i, "th host: ", host, " for reason: ", err, " .Skipping...")
			}
			defer conn.Close()
			handshake := make([]byte, 68)
			handshake[0] = 19
			copy(handshake[1:20], []byte("BitTorrent Protocol"))
			copy(handshake[28:48], hash[:])
			copy(handshake[48:68], peerid[:])
			conn.Write(handshake)

			resp := make([]byte, 68)
			io.ReadFull(conn, resp)
			if string(resp[28:48]) != string(hash[:]) {
				conn.Close()
				continue
			}

		}

	}

	if len(peers.peers_v6) != 0 {
		for i, v := range peers.peers_v6 {
			WRN(i, v)
		}

	}

}
