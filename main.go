package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	lp := flag.Uint("listenPort", 0, "Listen port")
	cp := flag.Uint("connectPort", 0, "Connect port")
	isServer := flag.Bool("server", false, "")

	flag.Parse()

	if *lp == 0 || *cp == 0 {
		fmt.Println("ポートが指定されていない")
		os.Exit(1)
	}

	if *lp != *cp {
		fmt.Println("listenPort と connectPort が同じでない (今後直す)")
		os.Exit(1)
	}

	if *isServer {
		go StartServer(int(*lp))
	} else {
		go StartClient(int(*cp))
	}

	<-make(chan struct{})
}
