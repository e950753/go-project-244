// Package diff строит дерево различий между двумя структурами данных, форматирование в него не входит
package diff

import (
	"maps"
	"reflect"
	"slices"
)

// Status описывает, как изменился ключ между двумя структурами
type Status string

const (
	// Added означает, что ключ есть только во второй структуре
	Added Status = "added"
	// Removed означает, что ключ есть только в первой структуре
	Removed Status = "removed"
	// Unchanged означает, что ключ есть в обеих структурах, значения равны
	Unchanged Status = "unchanged"
	// Updated означает, что ключ есть в обеих структурах, значения различаются
	Updated Status = "updated"
	// Nested означает, что в обеих структурах по ключу лежит объект, различия в Children
	Nested Status = "nested"
)

// Node описывает различие для одного ключа
type Node struct {
	// Key - имя ключа
	Key string
	// Status - вид изменения
	Status Status
	// OldValue - значение из первой структуры (Removed, Unchanged, Updated)
	OldValue any
	// NewValue - значение из второй структуры (Added, Unchanged, Updated)
	NewValue any
	// Children - различия вложенных объектов (только для Nested)
	Children []Node
}

// Build вычисляет различия между a и b, узлы отсортированы по ключу,
// дети строятся только если в обеих структурах по ключу лежит объект
func Build(a, b map[string]any) []Node {
	keys := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	nodes := make([]Node, 0, len(keys))
	for _, k := range slices.Sorted(maps.Keys(keys)) {
		av, inA := a[k]
		bv, inB := b[k]

		switch {
		case !inA:
			nodes = append(nodes, Node{Key: k, Status: Added, NewValue: bv})
		case !inB:
			nodes = append(nodes, Node{Key: k, Status: Removed, OldValue: av})
		default:
			am, aIsMap := av.(map[string]any)
			bm, bIsMap := bv.(map[string]any)
			switch {
			case aIsMap && bIsMap:
				nodes = append(nodes, Node{Key: k, Status: Nested, Children: Build(am, bm)})
			case reflect.DeepEqual(av, bv):
				nodes = append(nodes, Node{Key: k, Status: Unchanged, OldValue: av, NewValue: bv})
			default:
				nodes = append(nodes, Node{Key: k, Status: Updated, OldValue: av, NewValue: bv})
			}
		}
	}
	return nodes
}
