package client

import (
	"fmt"
	"github.com/jfbramlett/go-dynamic-runner/factories"
	"github.com/jfbramlett/go-dynamic-runner/routeguide"
	"github.com/jfbramlett/go-dynamic-runner/runner"
	"google.golang.org/grpc"
	"log"
)

func RunClient() {
	// register the set of client we are using (this is done to make sure they are pulled in to the build artifact
	factories.GlobalClientFactory.RegisterClient(routeguide.NewRouteGuideClient(newGrpcClient("localhost:2112")))

	// configure out test suite
	suite, err := runner.NewRunSuite("testdata/runsuite.json")
	if err != nil {
		log.Fatalln(err)
		return
	}

	// run the tests
	results := suite.Run()

	// this is just collecting/printing results
	passedCount := 0
	failedCount := 0
	for _, r := range results {
		if r.Passed {
			log.Println(fmt.Sprintf("Test %s - PASSED", r.Name))
			passedCount++
		} else {
			log.Println(fmt.Sprintf("Test %s - FAILED", r.Name))
			failedCount++
		}
	}

	log.Println(fmt.Sprintf("Total Tests: %d", passedCount + failedCount))
	log.Println(fmt.Sprintf("Passed Tests: %d", passedCount))
	log.Println(fmt.Sprintf("Failed Tests: %d", failedCount))


}


func  newGrpcClient(host string) *grpc.ClientConn {
	// prep our grpc env
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Fatalln("fail to dial ", err)
		return nil
	}

	return conn
}

