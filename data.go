package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func get_pieces(peers Peers, hash [20]byte, peerid [20]byte, cfg T_config) {
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
			if !bytes.Equal(resp[28:48], hash[:]) {
				conn.Close()
				continue
			}

			conn.Write([]byte{0, 0, 0, 1, 2})
			choked := true
			bitfield := []byte{}
			piece_cnt := cfg.n_pieces
			for {
				tgt := -1
				id, payload, err := read_data(conn)
				if err != nil {
					conn.Close()
					break
				}
				switch id {
				case 0:
					choked = true
				case 1:
					choked = false
				case 4:
					bitfield = append(bitfield, payload...)
				case 5:
					bitfield = payload

				case 255:

				}
				if choked {
					continue
				}

				for i := 0; i < piece_cnt; i++ {
					if contains_piece(bitfield, i) {
						tgt = i
						break
					}
				}

			}

		}

	}

	if len(peers.peers_v6) != 0 {
		for i, v := range peers.peers_v6 {
			WRN(i, v)
		}

	}

}

func contains_piece(buff []byte, idx int) bool {
	shift, offset := idx/8, 7-(idx%8)
	return (buff[shift]>>offset)&1 == 1

}

func read_data(conn net.Conn) (id byte, payload []byte, err error) {
	len_b := make([]byte, 4)
	if _, err = io.ReadFull(conn, len_b); err != nil {
		return
	}
	length := binary.BigEndian.Uint32(len_b)
	if length == 0 {
		return 255, nil, nil
	}
	msg := make([]byte, length)
	if _, err = io.ReadFull(conn, msg); err != nil {
		return
	}
	return msg[0], msg[1:], nil
}

func get_block(conn net.Conn, idx, strt, len int) {
	req := make()

}
