// Package main (example) shows additional usage patterns for the worker pool,
// including tasks that may return errors and processing results as they arrive.
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hasanulkabir-md/go-worker-pool/worker"
)

// divideTask divides 100 by the stored divisor.
// It returns an error when the divisor is zero to demonstrate error propagation.
type divideTask struct {
	divisor int
}

// Process implements task.Task.
func (d divideTask) Process() (interface{}, error) {
	if d.divisor == 0 {
		return nil, errors.New("division by zero")
	}
	return 100 / d.divisor, nil
}

func main() {
	// Use 4 workers.
	pool := worker.NewPool(4)
	pool.Start()

	// Submit 12 tasks, one of which will trigger an error.
	divisors := []int{1, 2, 4, 5, 0, 10, 20, 25, 50, 100, 3, 7}
	go func() {
		for i, d := range divisors {
			pool.Submit(i+1, divideTask{divisor: d})
		}
		pool.Stop()
	}()

	fmt.Println("Division results (100 / divisor):")
	for result := range pool.Results() {
		if result.Err != nil {
			log.Printf("  task %d error: %v\n", result.TaskID, result.Err)
			continue
		}
		fmt.Printf("  task %d: 100 / %d = %v\n", result.TaskID, divisors[result.TaskID-1], result.Value)
	}
}
