package async

import (
	"testing"
	"time"
)

type mockSink struct {
	logs [][]byte
}

func (m *mockSink) Write(data []byte) error {
	m.logs = append(m.logs, data)
	return nil
}

func TestWorkerBasic(t *testing.T) {
	worker := NewWorker(10)
	mock := &mockSink{}

	worker.Enqueue([]byte("log1"), []Sink{mock})
	worker.Enqueue([]byte("log2"), []Sink{mock})
	worker.Enqueue([]byte("log3"), []Sink{mock})

	worker.Stop()

	if len(mock.logs) != 3 {
		t.Errorf("expected 3 logs, got %d", len(mock.logs))
	}
}

func TestWorkerMultipleSinks(t *testing.T) {
	worker := NewWorker(10)
	mock1 := &mockSink{}
	mock2 := &mockSink{}

	worker.Enqueue([]byte("log1"), []Sink{mock1, mock2})

	worker.Stop()

	if len(mock1.logs) != 1 {
		t.Errorf("sink1: expected 1 log, got %d", len(mock1.logs))
	}

	if len(mock2.logs) != 1 {
		t.Errorf("sink2: expected 1 log, got %d", len(mock2.logs))
	}
}

func TestWorkerHighThroughput(t *testing.T) {
	worker := NewWorker(100)
	mock := &mockSink{}

	count := 1000
	for i := 0; i < count; i++ {
		worker.Enqueue([]byte("log"), []Sink{mock})
	}

	worker.Stop()

	if len(mock.logs) != count {
		t.Errorf("expected %d logs, got %d", count, len(mock.logs))
	}
}

func TestWorkerGracefulShutdown(t *testing.T) {
	worker := NewWorker(10)
	mock := &mockSink{}

	for i := 0; i < 5; i++ {
		worker.Enqueue([]byte("log"), []Sink{mock})
	}

	// Stop should wait for all logs to be processed
	done := make(chan bool)
	go func() {
		worker.Stop()
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("Stop() timed out")
	}

	if len(mock.logs) != 5 {
		t.Errorf("expected 5 logs after shutdown, got %d", len(mock.logs))
	}
}
