package client

import (
	"fmt"
	"github.com/jfbramlett/go-dynamic-runner/routeguide"
	"github.com/jfbramlett/go-dynamic-runner/runner"
	"google.golang.org/grpc"
	"log"
	"time"
)

func RunClient(configFile string) {
	// register the set of client we are using (this is done to make sure they are pulled in to the build artifact
	clients := map[string] interface{} {"routeguide.RouteGuideClient" : routeguide.NewRouteGuideClient(newGrpcClient("localhost:2112"))}

	// configure out test suite
	suite, err := runner.NewRunSuite(configFile, clients)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// run the tests
	results := suite.Run()

	// this is just collecting/printing results
	printResults(results)

}

func printResults(results runner.RunSuiteResult) {
	totalDuration := int64(results.Duration/time.Microsecond)
	log.Println(fmt.Sprintf("Test %d - RUN", results.TotalTests))
	log.Println(fmt.Sprintf("Test %d - PASSED", results.Passed))
	log.Println(fmt.Sprintf("Test %d - FAILED", results.Failed))

	avgMicro := totalDuration / int64(results.TotalTests)
	log.Println(fmt.Sprintf("Avg execution time: %d micros", avgMicro))
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

