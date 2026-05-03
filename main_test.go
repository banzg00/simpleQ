package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestConcurrentSubmitAndClaim(t *testing.T) {
	orch := &Orchestrator{Tasks: make(map[string]*Task)}

	// 5 producers submitting concurrently
	var submitWg sync.WaitGroup
	for i := range 5 {
		submitWg.Add(1)
		go func() {
			defer submitWg.Done()
			task := Task{
				ID:     fmt.Sprintf("task-%d", i),
				Status: StatusPending,
			}
			orch.Submit(task)
		}()
	}
	submitWg.Wait()

	// 3 workers claiming concurrently
	var claimWg sync.WaitGroup
	for i := range 3 {
		claimWg.Add(1)
		go func() {
			defer claimWg.Done()
			orch.Claim(fmt.Sprintf("w%d", i))
		}()
	}
	claimWg.Wait()

	pendningTasks := 0
	inProgressTasks := 0
	for _, task := range orch.Tasks {
		if task.Status == StatusPending {
			pendningTasks++
		}
		if task.Status == StatusInProgress {
			inProgressTasks++
		}
	}
	assert.Equal(t, pendningTasks, 2)
	assert.Equal(t, inProgressTasks, 3)
}
