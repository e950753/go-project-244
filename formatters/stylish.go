package formatters

import (
	"code/diff"
	"fmt"
	"maps"
	"slices"
	"strings"
)

const indentSize = 4

// отображает дерево различий в виде вложенных блоков с префиксами +, - и пробелом
func Stylish(nodes []diff.Node) string {
	return renderNodes(nodes, 1)
}

func renderNodes(nodes []diff.Node, depth int) string {
	signIndent := strings.Repeat(" ", depth*indentSize-2)
	lines := []string{"{"}

	for _, n := range nodes {
		switch n.Status {
		case diff.Added:
			lines = append(lines, line(signIndent, "+ ", n.Key, n.NewValue, depth))
		case diff.Removed:
			lines = append(lines, line(signIndent, "- ", n.Key, n.OldValue, depth))
		case diff.Unchanged:
			lines = append(lines, line(signIndent, "  ", n.Key, n.OldValue, depth))
		case diff.Updated:
			lines = append(lines,
				line(signIndent, "- ", n.Key, n.OldValue, depth),
				line(signIndent, "+ ", n.Key, n.NewValue, depth),
			)
		case diff.Nested:
			lines = append(lines,
				fmt.Sprintf("%s  %s: %s", signIndent, n.Key, renderNodes(n.Children, depth+1)))
		}
	}

	lines = append(lines, strings.Repeat(" ", (depth-1)*indentSize)+"}")
	return strings.Join(lines, "\n")
}

func line(indent, sign, key string, value any, depth int) string {
	return fmt.Sprintf("%s%s%s: %s", indent, sign, key, stringify(value, depth))
}

func stringify(value any, depth int) string {
	m, ok := value.(map[string]any)
	if !ok {
		return scalar(value)
	}

	innerIndent := strings.Repeat(" ", (depth+1)*indentSize)
	lines := []string{"{"}
	for _, k := range slices.Sorted(maps.Keys(m)) {
		lines = append(lines, fmt.Sprintf("%s%s: %s", innerIndent, k, stringify(m[k], depth+1)))
	}
	lines = append(lines, strings.Repeat(" ", depth*indentSize)+"}")
	return strings.Join(lines, "\n")
}
