package sync

import (
	"testing"
	"time"

	"github.com/zzjbattlefield/hajimi-go/internal/extractor"
	"github.com/zzjbattlefield/hajimi-go/internal/logger"
)

func TestNewSyncManager(t *testing.T) {
	logger := logger.New()

	sm := NewSyncManager(logger, 10, 30*time.Second, 3)

	if sm == nil {
		t.Error("Expected a sync manager, got nil")
	}

	if sm.batchSize != 10 {
		t.Errorf("Expected batch size to be 10, got %d", sm.batchSize)
	}

	if sm.flushInterval != 30*time.Second {
		t.Errorf("Expected flush interval to be 30 seconds, got %v", sm.flushInterval)
	}

	if sm.maxRetries != 3 {
		t.Errorf("Expected max retries to be 3, got %d", sm.maxRetries)
	}
}

func TestAddToQueue(t *testing.T) {
	logger := logger.New()
	sm := NewSyncManager(logger, 10, 30*time.Second, 3)

	secret := extractor.Secret{
		Type:   "test_key",
		Value:  "test_value",
		File:   "test.txt",
		Repo:   "test/repo",
		Line:   1,
		Commit: "abc123",
	}

	// Add a secret to the queue
	sm.AddToQueue(secret)

	// Check that the queue has one item
	if len(sm.queue) != 1 {
		t.Errorf("Expected queue length to be 1, got %d", len(sm.queue))
	}

	// Get the item from the queue
	item := <-sm.queue

	if item.Secret.Type != "test_key" {
		t.Errorf("Expected secret type to be 'test_key', got '%s'", item.Secret.Type)
	}

	if item.Secret.Value != "test_value" {
		t.Errorf("Expected secret value to be 'test_value', got '%s'", item.Secret.Value)
	}
}

func TestGetQueueLength(t *testing.T) {
	logger := logger.New()
	sm := NewSyncManager(logger, 10, 30*time.Second, 3)

	// Check initial queue length
	if sm.GetQueueLength() != 0 {
		t.Errorf("Expected initial queue length to be 0, got %d", sm.GetQueueLength())
	}

	// Add a secret to the queue
	secret := extractor.Secret{
		Type:   "test_key",
		Value:  "test_value",
		File:   "test.txt",
		Repo:   "test/repo",
		Line:   1,
		Commit: "abc123",
	}

	sm.AddToQueue(secret)

	// Check queue length after adding an item
	if sm.GetQueueLength() != 1 {
		t.Errorf("Expected queue length to be 1, got %d", sm.GetQueueLength())
	}
}

func TestGetBatchQueueLength(t *testing.T) {
	logger := logger.New()
	sm := NewSyncManager(logger, 10, 30*time.Second, 3)

	// Check initial batch queue length
	if sm.GetBatchQueueLength() != 0 {
		t.Errorf("Expected initial batch queue length to be 0, got %d", sm.GetBatchQueueLength())
	}

	// Send a batch to the batch queue
	batch := []*SyncItem{
		{
			Secret: extractor.Secret{
				Type:   "test_key",
				Value:  "test_value",
				File:   "test.txt",
				Repo:   "test/repo",
				Line:   1,
				Commit: "abc123",
			},
			Timestamp: time.Now(),
			Attempts:  0,
		},
	}

	sm.sendBatch(batch)

	// Check batch queue length after sending a batch
	if sm.GetBatchQueueLength() != 1 {
		t.Errorf("Expected batch queue length to be 1, got %d", sm.GetBatchQueueLength())
	}
}

func TestProcessBatch(t *testing.T) {
	logger := logger.New()
	sm := NewSyncManager(logger, 10, 30*time.Second, 3)

	// Create a batch of items
	batch := []*SyncItem{
		{
			Secret: extractor.Secret{
				Type:   "test_key1",
				Value:  "test_value1",
				File:   "test1.txt",
				Repo:   "test/repo",
				Line:   1,
				Commit: "abc123",
			},
			Timestamp: time.Now(),
			Attempts:  0,
		},
		{
			Secret: extractor.Secret{
				Type:   "test_key2",
				Value:  "test_value2",
				File:   "test2.txt",
				Repo:   "test/repo",
				Line:   2,
				Commit: "def456",
			},
			Timestamp: time.Now(),
			Attempts:  0,
		},
	}

	// Process the batch
	failedItems := sm.processBatch(batch)

	// In our implementation, no items should fail
	if len(failedItems) != 0 {
		t.Errorf("Expected no failed items, got %d", len(failedItems))
	}
}

func TestProcessBatchWithRetries(t *testing.T) {
	logger := logger.New()
	sm := NewSyncManager(logger, 10, 30*time.Second, 3)

	// Create a batch of items, including one that will fail on first attempt
	batch := []*SyncItem{
		{
			Secret: extractor.Secret{
				Type:   "github_token", // This will fail on first attempt in our implementation
				Value:  "test_value1",
				File:   "test1.txt",
				Repo:   "test/repo",
				Line:   1,
				Commit: "abc123",
			},
			Timestamp: time.Now(),
			Attempts:  0,
		},
		{
			Secret: extractor.Secret{
				Type:   "test_key2",
				Value:  "test_value2",
				File:   "test2.txt",
				Repo:   "test/repo",
				Line:   2,
				Commit: "def456",
			},
			Timestamp: time.Now(),
			Attempts:  0,
		},
	}

	// Process the batch
	failedItems := sm.processBatch(batch)

	// In our implementation, the github_token should fail
	if len(failedItems) != 1 {
		t.Errorf("Expected 1 failed item, got %d", len(failedItems))
	}

	if failedItems[0].Secret.Type != "github_token" {
		t.Errorf("Expected failed item to be 'github_token', got '%s'", failedItems[0].Secret.Type)
	}
}
