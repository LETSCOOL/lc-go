package executor

import (
	"context"
	"lc-go/lg"
	"math"
	"runtime"
	"sync"
)

type Microservice struct {
	numOfRoutines int
}

func (m *Microservice) Run() {
	max := lg.Ife(m.numOfRoutines <= 0, int(math.Min(float64(runtime.NumCPU()), float64(runtime.GOMAXPROCS(0)))), m.numOfRoutines)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(max)
	for i := 0; i < max; i++ {
		go func() {
			context.TODO()
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}
