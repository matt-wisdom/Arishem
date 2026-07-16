package jobs

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

var (
	// Worker pool settings
	maxWorkers = 4
	maxPerOrg  = 2

	// Task queues with large buffers
	pentestQueue = make(chan *LLMPentestTask, 200)
	scanQueue    = make(chan *ScanTask, 200)

	// Concurrency tracking and Registry
	activeTasksMu sync.Mutex
	activeTasks   = make(map[string]context.CancelFunc)
	orgTaskCounts = make(map[string]int)
)

// StartWorkerPool initializes the background workers
func StartWorkerPool(ctx context.Context) {
	for i := 0; i < maxWorkers; i++ {
		go llmWorker(ctx)
		go scanWorker(ctx)
	}
	slog.Info("Worker pools started", slog.Int("workers", maxWorkers), slog.Int("max_per_org", maxPerOrg))
}

// RegisterActiveTask registers a task execution context cancel function, enforcing per-org rate limiting
func RegisterActiveTask(id string, orgID string, cancel context.CancelFunc) bool {
	activeTasksMu.Lock()
	defer activeTasksMu.Unlock()

	// Enforce resource limit per user/organization
	if orgTaskCounts[orgID] >= maxPerOrg {
		return false
	}

	activeTasks[id] = cancel
	orgTaskCounts[orgID]++
	return true
}

// DeregisterActiveTask removes the task and decrements organization concurrency count
func DeregisterActiveTask(id string, orgID string) {
	activeTasksMu.Lock()
	defer activeTasksMu.Unlock()

	if _, exists := activeTasks[id]; exists {
		delete(activeTasks, id)
		orgTaskCounts[orgID]--
		if orgTaskCounts[orgID] < 0 {
			orgTaskCounts[orgID] = 0
		}
	}
}

// CancelActiveTask cancels a running task context, instantly killing its underlying exec.Cmd
func CancelActiveTask(id string) bool {
	activeTasksMu.Lock()
	defer activeTasksMu.Unlock()

	if cancel, exists := activeTasks[id]; exists {
		cancel()
		delete(activeTasks, id)
		return true
	}
	return false
}

// llmWorker pulls and runs LLMPentestTasks
func llmWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-pentestQueue:
			runCtx, cancel := context.WithCancel(ctx)

			if !RegisterActiveTask(task.RunID, task.OrgID, cancel) {
				cancel()
				// Rate-limited: Wait 1 second and put it back in queue
				go func(t *LLMPentestTask) {
					time.Sleep(1 * time.Second)
					select {
					case pentestQueue <- t:
					default:
						go func() { pentestQueue <- t }()
					}
				}(task)
				continue
			}

			processLLMPentestTask(runCtx, task)
			DeregisterActiveTask(task.RunID, task.OrgID)
		}
	}
}

// scanWorker pulls and runs ScanTasks
func scanWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-scanQueue:
			runCtx, cancel := context.WithCancel(ctx)

			if !RegisterActiveTask(task.ScanID, task.OrgID, cancel) {
				cancel()
				// Rate-limited: Wait 1 second and put it back in queue
				go func(t *ScanTask) {
					time.Sleep(1 * time.Second)
					select {
					case scanQueue <- t:
					default:
						go func() { scanQueue <- t }()
					}
				}(task)
				continue
			}

			processScanTask(runCtx, task)
			DeregisterActiveTask(task.ScanID, task.OrgID)
		}
	}
}
