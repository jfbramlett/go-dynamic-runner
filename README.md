# Description
Library that can be used to dynamically execute one or more Go methods with a focus on running GRPC.


# Run the sample code
1. Generate the proto using the make_proto.sh
2. Run example.go from cmd/example

The main takes 1 argument, --server to run the server and --client to run the client

The client will use the configuration defined in client/testdata/runsuite.json to run a set of methods against the server.