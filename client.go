package main

import (
	"context"
	"log"
	"net"
	"time"
)

type clientCon struct {
	Con    *net.TCPConn
	Cancel context.CancelFunc
}

var clientConMap = make(map[string]clientCon)

func StartClient(port int) {
	ucon, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		log.Println("Failed to listen UDP " + err.Error())
		return
	}

	readTCP := func(ctx context.Context, tcon *net.TCPConn, ucon *net.UDPConn, remote *net.UDPAddr) {
		buf := make([]byte, 4*1024)
		for {
			select {
			case <-ctx.Done():
				log.Println("cancel")
				return
			default:
			}

			n, err := tcon.Read(buf)
			if err != nil {
				log.Println("Failed to read TCP " + err.Error())
				tcon.Close()
				delete(clientConMap, remote.String())
				return
			}
			if _, err := ucon.WriteToUDP(buf[:n], remote); err != nil {
				log.Println("Failed to write UDP " + err.Error())
				continue
			}
		}
	}

	readUDP := func(ucon *net.UDPConn) {
		lbuf := make([]byte, 4*1024)
		for {
			n, remote, err := ucon.ReadFromUDP(lbuf)
			if err != nil {
				log.Println("Failed to read UDP " + err.Error())
				continue
			}

			cc, ok := clientConMap[remote.String()]
			if !ok {
				c, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
				if err != nil {
					log.Println("Failed to connect TCP " + err.Error())
					continue
				}
				c.SetDeadline(time.Now().Add(30 * time.Minute))

				ctx, cancel := context.WithCancel(context.Background())

				go readTCP(ctx, c, ucon, remote)

				clientConMap[remote.String()] = clientCon{
					Con:    c,
					Cancel: cancel,
				}
				cc = clientCon{
					Con:    c,
					Cancel: cancel,
				}
				log.Println("Connect TCP")
			}

			if _, err := cc.Con.Write(lbuf[:n]); err != nil {
				log.Println("Failed to write TCP: " + err.Error())
				continue
			}
		}
	}

	go readUDP(ucon)
}
