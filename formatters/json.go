package formatters

import (
	"code/diff"
	"encoding/json"
	"fmt"
)

func JSON(nodes []diff.Node) (string, error) {
	out, err := json.MarshalIndent(toJSONNodes(nodes), "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	return string(out), nil
}

func toJSONNodes(nodes []diff.Node) []map[string]any {
	result := make([]map[string]any, 0, len(nodes))
	for _, n := range nodes {
		m := map[string]any{
			"key":    n.Key,
			"status": string(n.Status),
		}
		switch n.Status {
		case diff.Added:
			m["value"] = n.NewValue
		case diff.Removed, diff.Unchanged:
			m["value"] = n.OldValue
		case diff.Updated:
			m["oldValue"] = n.OldValue
			m["newValue"] = n.NewValue
		case diff.Nested:
			m["children"] = toJSONNodes(n.Children)
		}
		result = append(result, m)
	}
	return result
}
