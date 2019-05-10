package client

import (
	"context"
	"github.com/jfbramlett/go-dynamic-runner/routeguide"
	"github.com/jfbramlett/go-dynamic-runner/runner"
	"time"
)

func RunClientDirect(loops int) {
	client := routeguide.NewRouteGuideClient(newGrpcClient("localhost:2112"))

	runResults := runner.RunSuiteResult{}
	for i := 0; i < loops; i++ {
		start := time.Now()
		runResults.TotalTests++

		proto.
		request := &routeguide.RouteRequest{Destination: "UNC", Email: "jon.snow@got.com", Userid: 10, TripHome: "yes", Uuid: "abc123",
			Pagination: &routeguide.Pagination{Size: 10, PageNum: 50}}

		_, err := client.FindRoute(context.Background(), request)

		if err == nil {
			runResults.Passed++
		} else {
			runResults.Failed++
		}
		runResults.Duration = runResults.Duration + time.Since(start)
	}

	// this is just collecting/printing results
	printResults(runResults)
}