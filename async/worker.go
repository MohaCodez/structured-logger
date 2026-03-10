package async

import (
	"fmt"
	"os"
	"sync/atomic"
)

type LogEntry struct {
	Data  []byte
	Sinks []Sink
}

type Sink interface {
	Write(data []byte) error
}

type Worker struct {
	queue        chan LogEntry
	done         chan struct{}
	dropOnFull   bool
	droppedCount uint64
	errorHandler func(error)
}

func NewWorker(bufferSize int, dropOnFull bool, errorHandler func(error)) *Worker {
	if errorHandler == nil {
		errorHandler = func(err error) {
			fmt.Fprintf(os.Stderr, "async worker: failed to write to sink: %v\n", err)
		}
	}
	
	w := &Worker{
		queue:        make(chan LogEntry, bufferSize),
		done:         make(chan struct{}),
		dropOnFull:   dropOnFull,
		errorHandler: errorHandler,
	}
	go w.run()
	return w
}

func (w *Worker) run() {
	for entry := range w.queue {
		for _, sink := range entry.Sinks {
			if err := sink.Write(entry.Data); err != nil {
				w.errorHandler(err)
			}
		}
	}
	close(w.done)
}

func (w *Worker) Enqueue(data []byte, sinks []Sink) {
	entry := LogEntry{Data: data, Sinks: sinks}
	
	if w.dropOnFull {
		// Non-blocking send: drop if buffer is full
		select {
		case w.queue <- entry:
		default:
			atomic.AddUint64(&w.droppedCount, 1)
		}
	} else {
		// Blocking send: wait for buffer space (backpressure)
		w.queue <- entry
	}
}

func (w *Worker) DroppedCount() uint64 {
	return atomic.LoadUint64(&w.droppedCount)
}

func (w *Worker) Stop() {
	close(w.queue)
	<-w.done
}
