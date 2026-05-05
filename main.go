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

// squareResult carries both the original number and its square so that the
// collector does not need to re-derive the input from the task ID.
type squareResult struct {
	input  int
	output int
}

// Process implements task.Task.  It squares the stored number and returns the result.
func (s squareTask) Process() (interface{}, error) {
	return squareResult{input: s.number, output: s.number * s.number}, nil
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
		sr := result.Value.(squareResult)
		fmt.Printf("  task %d: %d² = %d\n", result.TaskID, sr.input, sr.output)
	}

	fmt.Println("\nAll tasks completed.")
}
