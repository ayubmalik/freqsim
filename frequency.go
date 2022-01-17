package freqsim

import (
	"time"

	"golang.org/x/exp/rand"
)

type Frequency float32

type RandomFrequencyMeter struct {
	Interval time.Duration
	f        Frequency
	done     chan bool
}

func (m *RandomFrequencyMeter) Read() Frequency {
	return m.f
}

func (m *RandomFrequencyMeter) Start() {
	m.done = make(chan bool)
	ticker := time.NewTicker(m.Interval * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.f = Frequency(50.0 + rand.NormFloat64()*0.025)
			case <-m.done:
				return
			}
		}
	}()
}

func (m *RandomFrequencyMeter) Stop() {
	if m.done != nil {
		m.done <- true
	}
}
