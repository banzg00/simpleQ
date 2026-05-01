package main

import (
	"encoding/json"
	"fmt"
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
	tasks map[string]*Task
	mu    sync.Mutex
}

func (o *Orchestrator) Submit(task Task) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	_, exists := o.tasks[task.ID]
	if exists {
		return fmt.Errorf("Task with ID: [%s] already exists in the queue", task.ID)
	}
	o.tasks[task.ID] = &task
	return nil
}

func (o *Orchestrator) Claim(workerID string) (*Task, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.tasks) == 0 {
		return nil, fmt.Errorf("Queue is empty")
	}

	for _, task := range o.tasks {
		if task.Status == StatusPending {
			task.WorkerID = &workerID
			task.Status = StatusInProgress
			task.UpdatedAt = time.Now()
			return task, nil
		}
	}

	return nil, fmt.Errorf("No tasks to be claimed")
}

func (o *Orchestrator) Complete(taskID string, status TaskStatus) error {
	if status != StatusFailed &&
		status != StatusDead &&
		status != StatusDone {
		return fmt.Errorf("Wrong terminal status")
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	task, exists := o.tasks[taskID]
	if !exists {
		return fmt.Errorf("Task with ID: [%s] does not exists in the queue", taskID)
	}
	task.Status = status
	task.UpdatedAt = time.Now()

	return nil
}

type Producer struct {
	submitter TaskSubmitter
}

type Worker struct {
	ID        string
	claimer   TaskClaimer
	completer TaskCompleter
}

func main() {
	orch := &Orchestrator{tasks: make(map[string]*Task)}
	prod := Producer{submitter: orch}
	worker := Worker{ID: "w1", claimer: orch, completer: orch}

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

	err := prod.submitter.Submit(t1)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}

	recievedTask, err := worker.claimer.Claim(worker.ID)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}
	payload := map[string]any{}
	err = json.Unmarshal(recievedTask.Payload, &payload)
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	}
	fmt.Println(payload)

	worker.completer.Complete(recievedTask.ID, StatusDone)
	fmt.Println("Done!")
}
