package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/ayubmalik/freqsim/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println(conn.GetState().String())
	client := pb.NewFrequencySimulatorClient(conn)

	ctx, empty := context.Background(), &emptypb.Empty{}
	frequency, err := client.Get(context.Background(), empty)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	client.Read(ctx, empty)
	log.Printf("F = %v", frequency.Value)

	stream, err := client.Read(ctx, empty)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}

	waitc := make(chan struct{})
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		log.Printf("Got F %f, %v\n", in.Value, in.Time)
	}

}
