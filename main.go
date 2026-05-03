package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusFailed     TaskStatus = "failed"
	StatusDone       TaskStatus = "done"
	StatusDead       TaskStatus = "dead"
)

type Task struct {
	ID         string
	Payload    []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
	WorkerID   *string
	TaskType   string
	RetryCount int
	MaxRetries *int
	ExpiresAt  time.Time
	Status     TaskStatus
}

type TaskSubmitter interface {
	Submit(task Task) error
}

type TaskClaimer interface {
	Claim(workerID string) (*Task, error)
}

// Used for completing a task by moving it to one of terminal states
type TaskCompleter interface {
	Complete(taskID string, status TaskStatus) error
}

type Orchestrator struct {
	Tasks      map[string]*Task
	MaxRetries int
	mu         sync.Mutex
}

func (o *Orchestrator) Submit(task Task) error {
	o.mu.Lock()

	_, exists := o.Tasks[task.ID]
	if exists {
		o.mu.Unlock()
		return fmt.Errorf("task with ID=%s already exists in the queue", task.ID)
	}
	o.Tasks[task.ID] = &task
	o.mu.Unlock()

	slog.Info("task submitted", "taskID", task.ID)
	return nil
}

func (o *Orchestrator) Claim(workerID string) (*Task, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.Tasks) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	for _, task := range o.Tasks {
		if time.Now().After(task.ExpiresAt) {
			continue
		}
		if task.Status == StatusPending {
			task.WorkerID = &workerID
			task.Status = StatusInProgress
			task.UpdatedAt = time.Now()
			return task, nil
		} else if task.Status == StatusFailed && task.RetryCount < o.getMaxTaskRetries(task) {
			task.WorkerID = &workerID
			task.Status = StatusInProgress
			task.RetryCount++
			task.UpdatedAt = time.Now()
			return task, nil
		}
	}

	return nil, fmt.Errorf("no tasks to be claimed")
}

func (o *Orchestrator) Complete(taskID string, status TaskStatus) error {
	if status != StatusFailed &&
		status != StatusDead &&
		status != StatusDone {
		return fmt.Errorf("wrong terminal status")
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	task, exists := o.Tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID=%s does not exists in the queue", taskID)
	}
	task.Status = status
	task.UpdatedAt = time.Now()

	return nil
}

func (o *Orchestrator) getMaxTaskRetries(task *Task) int {
	if task.MaxRetries != nil && *task.MaxRetries >= 0 {
		return *task.MaxRetries
	}
	return o.MaxRetries
}

type Producer struct {
	submitter TaskSubmitter
}

type Worker struct {
	ID        string
	claimer   TaskClaimer
	completer TaskCompleter
}

func (worker *Worker) Process() error {
	task, err := worker.claimer.Claim(worker.ID)
	if err != nil {
		return fmt.Errorf("claim failed: %w", err)
	}
	slog.Info("task claimed", "workerID", worker.ID, "taskID", task.ID)

	slog.Info("processing task", "taskID", task.ID)
	time.Sleep(3 * time.Second)
	payload := map[string]any{}
	err = json.Unmarshal(task.Payload, &payload)
	if err != nil {
		return fmt.Errorf("processing failed: %w\n", err)
	}
	slog.Info("task payload", "taskID", task.ID, "payload", payload)

	err = worker.completer.Complete(task.ID, StatusDone)
	if err != nil {
		return fmt.Errorf("complete failed: %w", err)
	}
	slog.Info("task completed", "workerID", worker.ID, "taskID", task.ID)
	return nil
}

func main() {
	orch := &Orchestrator{Tasks: make(map[string]*Task), MaxRetries: 3}
	prod := Producer{submitter: orch}

	// 1. A freshly submitted task, waiting to be picked up
	t1 := Task{
		ID:         "task-001",
		TaskType:   "llm.summarize",
		Payload:    []byte(`{"model":"gpt-4","prompt":"Summarize this document..."}`),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Status:     StatusPending,
		RetryCount: 0,
		WorkerID:   nil, // no worker yet
	}
	t2 := Task{
		ID:         "task-002",
		TaskType:   "llm.summarize",
		Payload:    []byte(`{"model":"gpt-4","prompt":"Summarize this document..."}`),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Status:     StatusPending,
		RetryCount: 0,
		WorkerID:   nil, // no worker yet
	}
	t3 := Task{
		ID:         "task-003",
		TaskType:   "llm.summarize",
		Payload:    []byte(`{"model":"gpt-4","prompt":"Summarize this document..."}`),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Status:     StatusPending,
		RetryCount: 0,
		WorkerID:   nil, // no worker yet
	}
	t4 := Task{
		ID:         "task-004",
		TaskType:   "llm.summarize",
		Payload:    []byte(`{"model":"gpt-4","prompt":"Summarize this document..."}`),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Status:     StatusPending,
		RetryCount: 0,
		WorkerID:   nil, // no worker yet
	}
	t5 := Task{
		ID:         "task-005",
		TaskType:   "llm.summarize",
		Payload:    []byte(`{"model":"gpt-4","prompt":"Summarize this document..."}`),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Status:     StatusPending,
		RetryCount: 0,
		WorkerID:   nil, // no worker yet
	}

	err := prod.submitter.Submit(t1)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}
	err = prod.submitter.Submit(t2)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}
	err = prod.submitter.Submit(t3)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}
	err = prod.submitter.Submit(t4)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}
	err = prod.submitter.Submit(t5)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}

	var wg sync.WaitGroup

	for i := range 3 {
		wg.Go(func() {
			id := fmt.Sprintf("w%d", i)
			worker := Worker{ID: id, claimer: orch, completer: orch}
			err := worker.Process()
			if err != nil {
				slog.Error("worker failed", "workerID", id, "error", err)
			}
		})
	}

	wg.Wait()

}
