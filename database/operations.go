package database

import (
	"fmt"
)

func (t *Table) InsertRow(row Row) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.validateRow(row); err != nil {
		return err
	}

	t.Rows = append(t.Rows, row)
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
			// Validate new values
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
