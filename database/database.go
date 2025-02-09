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
	Name   string
	Tables map[string]*Table
	mutex  sync.RWMutex
}

type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
	mutex   sync.RWMutex
}

type Column struct {
	Name     string
	DataType string
}

type Row map[string]interface{}

func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Tables: make(map[string]*Table),
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
