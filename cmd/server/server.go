package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ayubmalik/freqsim"
	"github.com/ayubmalik/freqsim/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type frequencySimulatorServer struct {
	pb.UnimplementedFrequencySimulatorServer
	meter *freqsim.RandomFrequencyMeter
}

func (s *frequencySimulatorServer) Get(context.Context, *empty.Empty) (*pb.Frequency, error) {
	freq := &pb.Frequency{Value: float64(s.meter.Read()), Time: timestamppb.Now()}
	return freq, nil
}

func (s *frequencySimulatorServer) Read(cfg *pb.Config, stream pb.FrequencySimulator_ReadServer) error {
	millis := cfg.IntervalMillis
	for {
		freq := &pb.Frequency{Value: float64(s.meter.Read()), Time: timestamppb.Now()}
		if err := stream.Send(freq); err != nil {
			return err
		}
		time.Sleep(time.Duration(millis) * time.Millisecond)
	}
}

func newRPCServer(rfm *freqsim.RandomFrequencyMeter) *frequencySimulatorServer {
	s := &frequencySimulatorServer{meter: rfm}
	return s
}

func startRPCServer(rfm *freqsim.RandomFrequencyMeter) (*grpc.Server, error) {
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	pb.RegisterFrequencySimulatorServer(grpcServer, newRPCServer(rfm))
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	log.Printf("started rpc server on port: %d\n", port)
	return grpcServer, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// TODO(ayubm) use embedded resource
	serverCert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
