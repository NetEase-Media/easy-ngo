package pb

//go:generate protoc -I. --go_out=paths=source_relative:../testdata --go-grpc_out=paths=source_relative:../testdata helloworld.proto
