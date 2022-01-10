package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ayubmalik/freqsim/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type frequencySimulatorServer struct {
	protobuf.UnimplementedFrequencySimulatorServer
}

func (*frequencySimulatorServer) Get(context.Context, *empty.Empty) (*protobuf.Frequency, error) {
	return &protobuf.Frequency{Value: 666.66}, nil
}

func newServer() *frequencySimulatorServer {
	s := &frequencySimulatorServer{}
	return s
}

func startRPCServer() (*grpc.Server, error) {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	protobuf.RegisterFrequencySimulatorServer(grpcServer, newServer())
	log.Printf("starting rpc on port: %d\n", port)
	go func() { grpcServer.Serve(lis) }()
	log.Printf("started")
	return grpcServer, nil
}
