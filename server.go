package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ayubmalik/freqsim/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type frequencySimulatorServer struct {
	pb.UnimplementedFrequencySimulatorServer
}

func (*frequencySimulatorServer) Get(context.Context, *empty.Empty) (*pb.Frequency, error) {
	return &pb.Frequency{Value: 666.66}, nil
}

func (*frequencySimulatorServer) Read(_ *empty.Empty, stream pb.FrequencySimulator_ReadServer) error {
	for i := 0; i < 10; i++ {
		freq := pb.Frequency{Value: 123.00, Time: timestamppb.Now()}
		if err := stream.Send(&freq); err != nil {
			return err
		}
	}
	return nil
}

func newServer() *frequencySimulatorServer {
	s := &frequencySimulatorServer{}
	return s
}

func startRPCServer() (*grpc.Server, error) {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFrequencySimulatorServer(grpcServer, newServer())
	log.Printf("starting rpc on port: %d\n", port)
	go func() { grpcServer.Serve(lis) }()
	log.Printf("started")
	return grpcServer, nil
}
