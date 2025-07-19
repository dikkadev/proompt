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

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into StringSlice", value)
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

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONMap", value)
	}

	return json.Unmarshal(bytes, m)
}
