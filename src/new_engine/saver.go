package new_engine

import (
	"sync/atomic"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/timer"
)

// Saver handles asynchronous I/O queueing and flushing to disk.
type Saver struct {
	saveQue chan func()
	queLen  atom // Number of tasks in the queue
}

// NewSaver initializes a Saver with a specified queue capacity.
func NewSaver(queueCapacity int) *Saver {
	s := &Saver{
		saveQue: make(chan func(), queueCapacity),
	}
	go s.run()
	return s
}

// Atomic int
type atom int32

func (a *atom) Add(v int32) {
	atomic.AddInt32((*int32)(a), v)
}

func (a *atom) Load() int32 {
	return atomic.LoadInt32((*int32)(a))
}

// QueueOutput queues a function for asynchronous execution.
func (s *Saver) QueueOutput(f func()) {
	if cuda.Synchronous {
		timer.Start("io")
	}
	s.queLen.Add(1)
	s.saveQue <- f
	if cuda.Synchronous {
		timer.Stop("io")
	}
}

// run continuously executes tasks from the saveQue channel.
func (s *Saver) run() {
	for f := range s.saveQue {
		f()
		s.queLen.Add(-1)
	}
}

// Drain waits until all queued tasks are executed.
func (s *Saver) Drain() {
	for s.queLen.Load() > 0 {
		select {
		default:
			time.Sleep(1 * time.Millisecond) // Wait for the last job to finish
		case f := <-s.saveQue:
			f()
			s.queLen.Add(-1)
		}
	}
}
