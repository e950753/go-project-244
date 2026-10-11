package diff

import (
	"maps"
	"reflect"
	"slices"
)

type Status string

const (
	Added     Status = "added"
	Removed   Status = "removed"
	Unchanged Status = "unchanged"
	Updated   Status = "updated"
	Nested    Status = "nested"
)

type Node struct {
	Key      string
	Status   Status
	OldValue any
	NewValue any
	Children []Node
}

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
