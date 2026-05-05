// main.go demonstrates the worker pool by concurrently squaring a list of numbers.
package main

import (
	"fmt"
	"log"

	"github.com/hasanulkabir-md/go-worker-pool/worker"
)

// squareTask is a simple task that squares an integer.
type squareTask struct {
	number int
}

// Process implements task.Task.  It squares the stored number and returns the result.
func (s squareTask) Process() (interface{}, error) {
	result := s.number * s.number
	return result, nil
}

func main() {
	const numWorkers = 3
	const numTasks = 10

	fmt.Printf("Starting worker pool with %d workers to process %d tasks...\n\n", numWorkers, numTasks)

	// Create and start the pool.
	pool := worker.NewPool(numWorkers)
	pool.Start()

	// Submit tasks: square numbers 1 through numTasks.
	go func() {
		for i := 1; i <= numTasks; i++ {
			pool.Submit(i, squareTask{number: i})
		}
		// Signal that no more tasks will be submitted.
		pool.Stop()
	}()

	// Collect and print results as they arrive.
	fmt.Println("Results:")
	for result := range pool.Results() {
		if result.Err != nil {
			log.Printf("task %d failed: %v\n", result.TaskID, result.Err)
			continue
		}
		fmt.Printf("  task %d: %v² = %v\n", result.TaskID, result.TaskID, result.Value)
	}

	fmt.Println("\nAll tasks completed.")
}
