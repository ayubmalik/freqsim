package main

import (
	"context"
	"log"
	"os"

	"github.com/ayubmalik/freqsim/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:8081", opts...)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println(conn.GetState().String())
	client := protobuf.NewFrequencySimulatorClient(conn)
	frequency, err := client.Get(context.Background(), &empty.Empty{})

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Printf("F = %v", frequency.Value)

}
