package main

import (
	"context"
	"log"
	"net"
)

type serverCon struct {
	Con    *net.UDPConn
	Cancel context.CancelFunc
}

var serverConMap = make(map[string]serverCon)

func StartServer(port int) {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		log.Println("Failed to listen TCP " + err.Error())
		return
	}

	readUDP := func(ctx context.Context, tcon *net.TCPConn, ucon *net.UDPConn) {
		log.Println("spawn")
		buf := make([]byte, 4*1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n, err := ucon.Read(buf)
			if err != nil {
				log.Println("Failed to read UDP " + err.Error())
				continue
			}
			if _, err := tcon.Write(buf[:n]); err != nil {
				log.Println("Failed to write TCP ", err.Error())
				continue
			}
		}
	}

	readTCP := func(tcon *net.TCPConn) {
		buf := make([]byte, 4*1024)
		for {
			sc, ok := serverConMap[tcon.RemoteAddr().String()]
			if !ok {
				c, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: port})
				if err != nil {
					log.Println("Failed to dial UDP " + err.Error())
					continue
				}
				ctx, cancel := context.WithCancel(context.Background())

				serverConMap[tcon.RemoteAddr().String()] = serverCon{
					Con:    c,
					Cancel: cancel,
				}
				go readUDP(ctx, tcon, c)
				sc = serverCon{
					Con:    c,
					Cancel: cancel,
				}
				log.Println("Connect UDP")
			}

			n, err := tcon.Read(buf)
			if err != nil {
				log.Println("Disconnect due to failed to read TCP " + err.Error())
				sc.Cancel()
				tcon.Close()
				sc.Con.Close()
				delete(serverConMap, tcon.RemoteAddr().String())
				return
			}

			if _, err := sc.Con.Write(buf[:n]); err != nil {
				log.Println("Failed to write UDP ", err.Error())
				continue
			}
		}
	}

	for {
		tcon, err := ln.AcceptTCP()
		if err != nil {
			log.Println("Failed to accept TCP ", err.Error())
			continue
		}
		go readTCP(tcon)
	}
}
