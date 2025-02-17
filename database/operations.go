package database

import (
	"fmt"
)

// Table struct is already declared in database.go

func (t *Table) InsertRow(row Row) error {
	if err := t.insertRow(row); err != nil {
		return err
	}
	return t.db.save()
}

func (t *Table) insertRow(row Row) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.validateRow(row); err != nil {
		return err
	}

	t.Rows = append(t.Rows, row)
	t.db.history.AddOperation(Operation{ // Use t.db instead of db
		Type:      "insert",
		TableName: t.Name,
		Data:      row,
	})
	return nil
}

func (t *Table) rollbackInsert() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Remove the last inserted row
	if len(t.Rows) > 0 {
		t.Rows = t.Rows[:len(t.Rows)-1]
	}
	return nil
}

func (t *Table) rollbackUpdate(newData, oldData Row) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i, row := range t.Rows {
		if matchConditions(row, newData) {
			t.Rows[i] = oldData
			return nil
		}
	}
	return fmt.Errorf("row not found for rollback")
}

func (t *Table) rollbackDelete(oldData Row) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.Rows = append(t.Rows, oldData)
	return nil
}

func (t *Table) Select(conditions map[string]interface{}) ([]Row, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if len(conditions) == 0 {
		return t.Rows, nil
	}

	var result []Row
	for _, row := range t.Rows {
		if matchConditions(row, conditions) {
			result = append(result, row)
		}
	}
	return result, nil
}

func (t *Table) Delete(conditions map[string]interface{}) (int, error) {
	deleted, err := t.delete(conditions)
	if err != nil {
		return 0, err
	}
	if err := t.db.save(); err != nil {
		return 0, err
	}
	return deleted, nil
}

func (t *Table) delete(conditions map[string]interface{}) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	deleted := 0
	newRows := make([]Row, 0)

	for _, row := range t.Rows {
		if !matchConditions(row, conditions) {
			newRows = append(newRows, row)
		} else {
			deleted++
		}
	}

	t.Rows = newRows
	return deleted, nil
}

func (t *Table) Update(conditions map[string]interface{}, updates map[string]interface{}) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	updated := 0
	for i, row := range t.Rows {
		if matchConditions(row, conditions) {
			// Store old data for history
			oldData := make(Row)
			for k, v := range row {
				oldData[k] = v
			}

			// Create new row with updates
			newRow := make(Row)
			for k, v := range row {
				newRow[k] = v
			}
			for k, v := range updates {
				if err := t.validateColumn(k, v); err != nil {
					return 0, err
				}
				newRow[k] = v
			}

			t.Rows[i] = newRow
			updated++

			// Add to history
			t.db.history.AddOperation(Operation{
				Type:      "update",
				TableName: t.Name,
				Data:      newRow,
				OldData:   oldData,
			})
		}
	}

	// Save changes to disk
	if updated > 0 {
		if err := t.db.save(); err != nil {
			return 0, err
		}
	}

	return updated, nil
}

func (t *Table) validateRow(row Row) error {
	for _, col := range t.Columns {
		val, exists := row[col.Name]
		if !exists {
			return fmt.Errorf("missing value for column %s", col.Name)
		}

		if !validateDataType(val, col.DataType) {
			return fmt.Errorf("invalid data type for column %s", col.Name)
		}
	}
	return nil
}

func (t *Table) validateColumn(name string, value interface{}) error {
	for _, col := range t.Columns {
		if col.Name == name {
			if !validateDataType(value, col.DataType) {
				return fmt.Errorf("invalid data type for column %s", name)
			}
			return nil
		}
	}
	return fmt.Errorf("column %s does not exist", name)
}

// validateDataType function is declared in database.go

func matchConditions(row Row, conditions map[string]interface{}) bool {
	for key, value := range conditions {
		if rowValue, exists := row[key]; !exists || rowValue != value {
			return false
		}
	}
	return true
}
