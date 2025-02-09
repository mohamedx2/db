package database

import (
	"fmt"
	"strconv"
	"strings"
)

type QueryCondition struct {
	Column   string
	Operator string
	Value    interface{}
}

func ParseWhereClause(whereClause string) (map[string]interface{}, error) {
	if whereClause == "" {
		return nil, nil
	}

	conditions := make(map[string]interface{})
	parts := strings.Split(whereClause, "AND")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Simple equals parsing for now
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid condition: %s", part)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		conditions[key] = parseValue(value)
	}

	return conditions, nil
}

func parseValue(value string) interface{} {
	// Remove quotes if present
	value = strings.Trim(value, `"'`)

	// Try parsing as bool
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Try parsing as int
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	// Try parsing as float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Return as string by default
	return value
}
