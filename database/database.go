package database

import (
	"fmt"
	"sync"
	"time"

	"db/storage" // Change this line to use module path
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
	storage *storage.Storage
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

func NewDatabase(name string, dataDir string) (*Database, error) {
	storage, err := storage.NewStorage(dataDir)
	if err != nil {
		return nil, err
	}

	db := &Database{
		Name:    name,
		Tables:  make(map[string]*Table),
		history: NewHistory(),
		storage: storage,
	}

	// Load existing data
	if err := db.load(); err != nil {
		return nil, err
	}

	return db, nil
}

type persistedData struct {
	Tables  map[string]*Table `json:"tables"`
	History []Operation       `json:"history"`
}

func (db *Database) save() error {
	db.mutex.RLock()
	data := persistedData{
		Tables:  db.Tables,
		History: db.history.GetOperations(),
	}
	db.mutex.RUnlock()

	return db.storage.Save("database.json", data)
}

func (db *Database) load() error {
	var data persistedData
	if err := db.storage.Load("database.json", &data); err != nil {
		return err
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	if data.Tables != nil {
		db.Tables = data.Tables
		// Restore database reference in tables
		for _, table := range db.Tables {
			table.db = db
		}
	}

	if data.History != nil {
		for _, op := range data.History {
			db.history.AddOperation(op)
		}
	}

	return nil
}

func (db *Database) CreateTable(name string, columns []Column) error {
	if err := db.createTable(name, columns); err != nil {
		return err
	}
	return db.save()
}

func (db *Database) createTable(name string, columns []Column) error {
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
