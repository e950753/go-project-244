// разбирает содержимое файлов конфигурации в map[string]any
package parsers

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// выбирает парсер по расширению файла
func Parse(data []byte, ext string) (map[string]any, error) {
	var result map[string]any

	switch strings.ToLower(ext) {
	case ".json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
		for k, v := range result {
			result[k] = normalize(v)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %q", ext)
	}

	if result == nil {
		result = map[string]any{}
	}
	return result, nil
}

// normalize приводит целые числа из YAML к float64, как в JSON
func normalize(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			x[k] = normalize(val)
		}
		return x
	case []any:
		for i, val := range x {
			x[i] = normalize(val)
		}
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case uint64:
		return float64(x)
	default:
		return v
	}
}
