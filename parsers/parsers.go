package parsers

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte, ext string) (map[string]any, error) {
	var result map[string]any

	switch strings.ToLower(ext) {
	case ".json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
	case
		".yml",
		".yaml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %q", ext)
	}

	if result == nil {
		result = map[string]any{}
	}
	return result, nil
}
