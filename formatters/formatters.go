// Package formatters преобразует дерево различий в текст разных форматов
package formatters

import (
	"code/diff"
	"fmt"
)

// формат вывода, используемый, если формат не указан
const DefaultFormat = "stylish"

// отображает дерево различий в выбранном формате: stylish, plain или json
// пустая строка = DefaultFormat
// неизвестный формат возвращает ошибку
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
