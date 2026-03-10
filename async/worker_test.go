package async

import (
	"sync"
	"testing"
	"time"
)

type mockSink struct {
	mu   sync.Mutex
	logs [][]byte
}

func (m *mockSink) Write(data []byte) error {
	m.mu.Lock()
	m.logs = append(m.logs, data)
	m.mu.Unlock()
	return nil
}

func (m *mockSink) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.logs)
}

func TestWorkerBasic(t *testing.T) {
	worker := NewWorker(10, false, nil)
	mock := &mockSink{}

	worker.Enqueue([]byte("log1"), []Sink{mock})
	worker.Enqueue([]byte("log2"), []Sink{mock})
	worker.Enqueue([]byte("log3"), []Sink{mock})

	worker.Stop()

	if mock.Len() != 3 {
		t.Errorf("expected 3 logs, got %d", mock.Len())
	}
}

func TestWorkerMultipleSinks(t *testing.T) {
	worker := NewWorker(10, false, nil)
	mock1 := &mockSink{}
	mock2 := &mockSink{}

	worker.Enqueue([]byte("log1"), []Sink{mock1, mock2})

	worker.Stop()

	if mock1.Len() != 1 {
		t.Errorf("sink1: expected 1 log, got %d", mock1.Len())
	}

	if mock2.Len() != 1 {
		t.Errorf("sink2: expected 1 log, got %d", mock2.Len())
	}
}

func TestWorkerHighThroughput(t *testing.T) {
	worker := NewWorker(100, false, nil)
	mock := &mockSink{}

	count := 1000
	for i := 0; i < count; i++ {
		worker.Enqueue([]byte("log"), []Sink{mock})
	}

	worker.Stop()

	if mock.Len() != count {
		t.Errorf("expected %d logs, got %d", count, mock.Len())
	}
}

func TestWorkerGracefulShutdown(t *testing.T) {
	worker := NewWorker(10, false, nil)
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

	if mock.Len() != 5 {
		t.Errorf("expected 5 logs after shutdown, got %d", mock.Len())
	}
}

func TestAsyncBufferFullBlocks(t *testing.T) {
	mock := &mockSink{}
	worker := NewWorker(1, false, nil) // Small buffer, blocking policy
	defer worker.Stop()

	// Fill the buffer
	worker.Enqueue([]byte("log1"), []Sink{mock})
	
	// This should block briefly but eventually succeed
	done := make(chan bool)
	go func() {
		worker.Enqueue([]byte("log2"), []Sink{mock})
		done <- true
	}()

	// Give some time for processing
	select {
	case <-done:
		// Success - the blocking call completed
	case <-time.After(100 * time.Millisecond):
		t.Error("blocking enqueue should have completed")
	}
}

func TestAsyncBufferFullDrops(t *testing.T) {
	mock := &mockSink{}
	worker := NewWorker(1, true, nil) // Small buffer, drop policy
	defer worker.Stop()

	// Fill the buffer
	worker.Enqueue([]byte("log1"), []Sink{mock})
	
	// These should be dropped immediately
	for i := 0; i < 10; i++ {
		worker.Enqueue([]byte("dropped"), []Sink{mock})
	}

	// Check dropped count
	if worker.DroppedCount() == 0 {
		t.Error("expected some logs to be dropped")
	}

	// Give time for processing
	time.Sleep(10 * time.Millisecond)
	
	// Should have processed only the first log
	if mock.Len() > 1 {
		t.Errorf("expected at most 1 log processed, got %d", mock.Len())
	}
}
