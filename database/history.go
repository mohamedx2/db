package database

import (
	"sync"
	"time"
)

type Operation struct {
	Type      string // "insert", "update", "delete"
	TableName string
	Data      Row
	OldData   Row // For updates and deletes
	Timestamp time.Time
}

type History struct {
	operations []Operation
	mutex      sync.RWMutex
}

func NewHistory() *History {
	return &History{
		operations: make([]Operation, 0),
	}
}

func (h *History) AddOperation(op Operation) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	op.Timestamp = time.Now()
	h.operations = append(h.operations, op)
}

func (h *History) GetOperations() []Operation {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	ops := make([]Operation, len(h.operations))
	copy(ops, h.operations)
	return ops
}

func (h *History) Clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.operations = make([]Operation, 0)
}
