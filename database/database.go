package database

import (
	"fmt"
	"sync"
	"time"
)

const (
	TypeString    = "string"
	TypeInt       = "int"
	TypeBool      = "bool"
	TypeFloat     = "float"
	TypeTimestamp = "timestamp"
)

type Database struct {
	Name    string
	Tables  map[string]*Table
	history *History
	mutex   sync.RWMutex
}

type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
	db      *Database // Add reference to parent database
	mutex   sync.RWMutex
}

type Column struct {
	Name     string
	DataType string
}

type Row map[string]interface{}

func NewDatabase(name string) *Database {
	return &Database{
		Name:    name,
		Tables:  make(map[string]*Table),
		history: NewHistory(),
	}
}

func (db *Database) CreateTable(name string, columns []Column) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.Tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}

	db.Tables[name] = &Table{
		Name:    name,
		Columns: columns,
		Rows:    make([]Row, 0),
		db:      db, // Set the database reference
	}
	return nil
}

func (db *Database) GetTable(name string) (*Table, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if table, exists := db.Tables[name]; exists {
		return table, nil
	}
	return nil, fmt.Errorf("table %s not found", name)
}

func (db *Database) Rollback() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	ops := db.history.GetOperations()
	if len(ops) == 0 {
		return fmt.Errorf("no operations to rollback")
	}

	// Rollback last operation
	lastOp := ops[len(ops)-1]
	table, exists := db.Tables[lastOp.TableName]
	if !exists {
		return fmt.Errorf("table %s not found", lastOp.TableName)
	}

	switch lastOp.Type {
	case "insert":
		return table.rollbackInsert()
	case "update":
		return table.rollbackUpdate(lastOp.Data, lastOp.OldData)
	case "delete":
		return table.rollbackDelete(lastOp.OldData)
	default:
		return fmt.Errorf("unknown operation type: %s", lastOp.Type)
	}
}

func (db *Database) GetHistory() *History {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.history
}

func validateDataType(value interface{}, dataType string) bool {
	switch dataType {
	case TypeString:
		_, ok := value.(string)
		return ok
	case TypeInt:
		// Handle both int and float64 for integers
		switch v := value.(type) {
		case float64:
			// Check if it's a whole number
			return v == float64(int(v))
		case int:
			return true
		default:
			return false
		}
	case TypeBool:
		_, ok := value.(bool)
		return ok
	case TypeFloat:
		// Accept both float64 and int for float fields
		_, isFloat := value.(float64)
		_, isInt := value.(int)
		return isFloat || isInt
	case TypeTimestamp:
		_, ok := value.(time.Time)
		return ok
	default:
		return false
	}
}
