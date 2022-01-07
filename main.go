package main

import (
	"fmt"
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

func main() {
	m := randomFreqMeter{interval: 100}
	m.run()

	for {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(m.read())
	}
}
