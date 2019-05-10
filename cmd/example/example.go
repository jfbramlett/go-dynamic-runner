package main

import (
	"flag"
	"fmt"
	"github.com/jfbramlett/go-dynamic-runner/cmd/client"
	"github.com/jfbramlett/go-dynamic-runner/cmd/server"
	"testing"
)

func main() {
	clientPtr := flag.Bool("client", false, "run in client mode")
	clientBenchPtr := flag.Bool("client-bench", false, "run in client benchmark mode")
	clientDirectPtr := flag.Bool("client-direct", false, "run in client direct mode")
	clientDirectBenchPtr := flag.Bool("client-direct-bench", false, "run in client direct benchmark mode")
	serverPtr := flag.Bool("server", false, "run in server mode")

	flag.Parse()

	if *clientPtr {
		client.RunClient("testdata/singlerunsuite.json")
	} else if *clientBenchPtr {
		fmt.Println(testing.Benchmark(TestRunClient))
	} else if *clientDirectPtr {
		client.RunClientDirect(10000)
	} else if *clientDirectBenchPtr {
		fmt.Println(testing.Benchmark(TestRunClientDirect))
	} else if *serverPtr {
		server.RunServer()
	} else {
		fmt.Println("need to specify client or server mode")
		flag.Usage()
	}
}


func TestRunClient(t *testing.B) {
	client.RunClient("testdata/singlerunsuite.json")
}


func TestRunClientDirect(t *testing.B) {
	client.RunClientDirect(1)
}