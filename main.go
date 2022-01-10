package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/exp/rand"
)

type frequency float32

type randomFreqMeter struct {
	interval time.Duration
	f        frequency
	done     chan bool
}

func (m *randomFreqMeter) read() frequency {
	return m.f
}

func (m *randomFreqMeter) run() {
	m.done = make(chan bool)
	ticker := time.NewTicker(m.interval * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.f = frequency(50.0 + rand.NormFloat64()*0.05)
			case <-m.done:
				return
			}
		}
	}()
}

func handleCtrlC() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx, cancel
}

func main() {
	// ctx, cancel := handleCtrlC()
	// defer cancel()

	// m := randomFreqMeter{interval: 100}
	// m.run()

	// var wg sync.WaitGroup
	// wg.Add(1)

	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			fmt.Println("exit loop")
	// 			return
	// 		default:
	// 			time.Sleep(100 * time.Millisecond)
	// 			fmt.Println(m.read())
	// 		}
	// 	}
	// }()
	// wg.Wait()

	server, _ := startRPCServer()

	time.Sleep(3 * time.Second)
	server.GracefulStop()
	fmt.Println("tis done")

}
