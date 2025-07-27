package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PromptType string

const (
	PromptTypeSystem PromptType = "system"
	PromptTypeUser   PromptType = "user"
	PromptTypeImage  PromptType = "image"
	PromptTypeVideo  PromptType = "video"
)

func (pt PromptType) Valid() bool {
	switch pt {
	case PromptTypeSystem, PromptTypeUser, PromptTypeImage, PromptTypeVideo:
		return true
	}
	return false
}

// StringSlice handles JSON marshaling for []string in database
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	case *string:
		if v != nil {
			bytes = []byte(*v)
		} else {
			*s = nil
			return nil
		}
	default:
		// Debug: log the actual type we're receiving
		return fmt.Errorf("cannot scan %T (value: %v) into StringSlice", value, value)
	}

	return json.Unmarshal(bytes, s)
}

// JSONMap handles JSON marshaling for map[string]interface{} in database
type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	case *string:
		if v != nil {
			bytes = []byte(*v)
		} else {
			*m = nil
			return nil
		}
	default:
		// Debug: log the actual type we're receiving
		return fmt.Errorf("cannot scan %T (value: %v) into JSONMap", value, value)
	}

	return json.Unmarshal(bytes, m)
}
