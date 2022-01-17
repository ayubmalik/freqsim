package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ayubmalik/freqsim"
	log "google.golang.org/grpc/grpclog"
)

func handleCtrlC() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx, cancel
}

func main() {
	log.SetLoggerV2(log.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 2))

	ctx, cancel := handleCtrlC()
	defer cancel()

	m := &freqsim.RandomFrequencyMeter{Interval: 100}
	m.Start()

	s, err := startRPCServer(m)
	if err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()

	m.Stop()
	s.GracefulStop()
	log.Infoln("Exit OK")
}
