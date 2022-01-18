package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/ayubmalik/freqsim/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	log "google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	log.SetLoggerV2(log.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 2))
	var opts []grpc.DialOption

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	conn, err := grpc.Dial(":50051", opts...)
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
