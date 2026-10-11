package formatters

import (
	"code/diff"
	"fmt"
)

const DefaultFormat = "stylish"

func Format(nodes []diff.Node, format string) (string, error) {
	switch format {
	case "", DefaultFormat:
		return Stylish(nodes), nil
	case "plain":
		return Plain(nodes), nil
	case "json":
		return JSON(nodes)
	default:
		return "", fmt.Errorf("unknown output format: %q", format)
	}
}
