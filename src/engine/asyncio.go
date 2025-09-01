package engine

import (
	"sync/atomic"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/timer"
)

// Asynchronous I/O queue flushes data to disk while simulation keeps running.
// See save.go, autosave.go

var (
	saveQue chan func() // passes save requests to runSaver for asyc IO
	queLen  atom        // # tasks in queue
)

func init() {

	saveQue = make(chan func())
	go runSaver()
}

// Atomic int
type atom int32

func (a *atom) Add(v int32) {
	atomic.AddInt32((*int32)(a), v)
}

func (a *atom) Load() int32 {
	return atomic.LoadInt32((*int32)(a))
}

func queOutput(f func()) {
	if cuda.Synchronous {
		timer.Start("io")
	}
	queLen.Add(1)
	saveQue <- f
	if cuda.Synchronous {
		timer.Stop("io")
	}
}

// Continuously executes tasks the from SaveQue channel.
func runSaver() {
	for f := range saveQue {
		f()
		queLen.Add(-1)
	}
}

// Finalizer function called upon program exit.
// Waits until all asynchronous output has been saved.
func drainOutput() {
	if saveQue == nil {
		return
	}
	for queLen.Load() > 0 {
		select {
		default:
			time.Sleep(1 * time.Millisecond) // other goroutine has the last job, wait for it to finish
		case f := <-saveQue:
			f()
			queLen.Add(-1)
		}
	}
}
