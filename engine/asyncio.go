package engine

import (
	"time"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/timer"
	"github.com/MathieuMoalic/amumax/util"
)

// Asynchronous I/O queue flushes data to disk while simulation keeps running.
// See save.go, autosave.go

var (
	saveQue chan func() // passes save requests to runSaver for asyc IO
	queLen  util.Atom   // # tasks in queue
)

func init() {
	DeclFunc("Flush", drainOutput, "Flush all pending output to disk.")

	saveQue = make(chan func())
	go runSaver()
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
