package main

import (
	"io"
	"net"
)

func tryClose(con io.Closer) {
	defer recover()
	con.Close()
}

func startTCPEchoServer(port int) {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		panic(err)
	}

	for {
		con, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func(con net.Conn) {
			buf := make([]byte, 4*1024)

			for {
				n, err := con.Read(buf)
				if err != nil {
					continue
				}
				if _, err := con.Write(buf[:n]); err != nil {
					continue
				}
			}
		}(con)
	}
}

func startUDPEchoServer(port int) {
	ucon, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 4*1024)
	for {
		n, remote, err := ucon.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if _, err := ucon.WriteToUDP(buf[:n], remote); err != nil {
			panic(err)
		}
	}
}
