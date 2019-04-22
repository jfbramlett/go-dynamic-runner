package main

import (
	"flag"
	"fmt"
	"github.com/jfbramlett/go-dynamic-runner/cmd/client"
	"github.com/jfbramlett/go-dynamic-runner/cmd/server"
)

func main() {
	clientPtr := flag.Bool("client", false, "run in client mode")
	serverPtr := flag.Bool("server", false, "run in server mode")

	flag.Parse()

	if *clientPtr {
		client.RunClient()
	} else if *serverPtr {
		server.RunServer()
	} else {
		fmt.Println("need to specify client or server mode")
		flag.Usage()
	}
}
