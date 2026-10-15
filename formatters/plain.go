package formatters

import (
	"code/diff"
	"fmt"
	"strings"
)

// отображает дерево различий в виде списка строк вида
func Plain(nodes []diff.Node) string {
	return strings.Join(plainLines(nodes, ""), "\n")
}

func plainLines(nodes []diff.Node, parent string) []string {
	var lines []string
	for _, n := range nodes {
		path := n.Key
		if parent != "" {
			path = parent + "." + n.Key
		}

		switch n.Status {
		case diff.Added:
			lines = append(lines,
				fmt.Sprintf("Property '%s' was added with value: %s", path, plainValue(n.NewValue)))
		case diff.Removed:
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", path))
		case diff.Updated:
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
				path, plainValue(n.OldValue), plainValue(n.NewValue)))
		case diff.Nested:
			lines = append(lines, plainLines(n.Children, path)...)
		case diff.Unchanged:
		}
	}
	return lines
}

func plainValue(v any) string {
	switch x := v.(type) {
	case map[string]any, []any:
		return "[complex value]"
	case string:
		return "'" + x + "'"
	default:
		return scalar(x)
	}
}
