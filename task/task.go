// Package task defines the Task interface and Result type used by the worker pool.
package task

// Task represents a unit of work to be executed by a worker.
// Implementations should encapsulate all data required to perform the computation.
type Task interface {
	// Process executes the task and returns a result value or an error.
	Process() (interface{}, error)
}

// Result holds the outcome of a completed task along with any error that occurred.
type Result struct {
	// TaskID is a user-supplied identifier that correlates this result with its task.
	TaskID int
	// Value is the value returned by Task.Process on success.
	Value interface{}
	// Err is non-nil when the task failed.
	Err error
}
