package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/ayubmalik/freqsim/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	log "google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	log.SetLoggerV2(log.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 2))
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewFrequencySimulatorClient(conn)

	// single value
	frequency, err := client.Get(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("F = %v\n", frequency.Value)

	// stream frequency
	stream, err := client.Read(context.Background(), &pb.Config{IntervalMillis: 200})
	if err != nil {
		log.Fatalf("%v %v\n", client, err)
	}

	waitc := make(chan struct{})
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive frequency: %v", err)
		}
		log.Infof("Got F %f, %s\n", in.Value, in.Time.AsTime().Format(time.RFC3339Nano))
	}

	<-waitc
}
