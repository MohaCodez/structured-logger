package async

import (
	"fmt"
	"os"
)

type LogEntry struct {
	Data []byte
	Sinks []Sink
}

type Sink interface {
	Write(data []byte) error
}

type Worker struct {
	queue chan LogEntry
	done  chan struct{}
}

func NewWorker(bufferSize int) *Worker {
	w := &Worker{
		queue: make(chan LogEntry, bufferSize),
		done:  make(chan struct{}),
	}
	go w.run()
	return w
}

func (w *Worker) run() {
	for entry := range w.queue {
		for _, sink := range entry.Sinks {
			if err := sink.Write(entry.Data); err != nil {
				fmt.Fprintf(os.Stderr, "async worker: failed to write to sink: %v\n", err)
			}
		}
	}
	close(w.done)
}

func (w *Worker) Enqueue(data []byte, sinks []Sink) {
	w.queue <- LogEntry{Data: data, Sinks: sinks}
}

func (w *Worker) Stop() {
	close(w.queue)
	<-w.done
}
