package formatters

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/diff"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name  string
		nodes []diff.Node
		want  string
	}{
		{"empty", nil, `{}`},
		{
			"added",
			[]diff.Node{{Key: "k", Status: diff.Added, NewValue: float64(1)}},
			`{"k": {"status": "added", "value": 1}}`,
		},
		{
			"removed",
			[]diff.Node{{Key: "k", Status: diff.Removed, OldValue: "v"}},
			`{"k": {"status": "removed", "value": "v"}}`,
		},
		{
			"unchanged",
			[]diff.Node{{Key: "k", Status: diff.Unchanged, OldValue: "v", NewValue: "v"}},
			`{"k": {"status": "unchanged", "value": "v"}}`,
		},
		{
			"updated keeps null",
			[]diff.Node{{Key: "k", Status: diff.Updated, OldValue: true, NewValue: nil}},
			`{"k": {"status": "updated", "oldValue": true, "newValue": null}}`,
		},
		{
			"added null keeps value key",
			[]diff.Node{{Key: "k", Status: diff.Added, NewValue: nil}},
			`{"k": {"status": "added", "value": null}}`,
		},
		{
			"complex value",
			[]diff.Node{{Key: "k", Status: diff.Added, NewValue: map[string]any{"x": []any{1}}}},
			`{"k": {"status": "added", "value": {"x": [1]}}}`,
		},
		{
			"nested",
			[]diff.Node{{
				Key: "n", Status: diff.Nested,
				Children: []diff.Node{{Key: "a", Status: diff.Removed, OldValue: 1}},
			}},
			`{"n": {"status": "nested", "children": {"a": {"status": "removed", "value": 1}}}}`,
		},
		{
			"empty nested",
			[]diff.Node{{Key: "n", Status: diff.Nested, Children: []diff.Node{}}},
			`{"n": {"status": "nested", "children": {}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSON(tt.nodes)
			require.NoError(t, err)
			require.JSONEq(t, tt.want, got)
		})
	}
}

func TestJSONIsIndented(t *testing.T) {
	got, err := JSON([]diff.Node{{Key: "k", Status: diff.Added, NewValue: float64(1)}})
	require.NoError(t, err)
	require.Contains(t, got, "\n  \"k\"")
}

func TestJSONMarshalError(t *testing.T) {
	// канал нельзя сериализовать в JSON
	_, err := JSON([]diff.Node{{Key: "k", Status: diff.Added, NewValue: make(chan int)}})
	require.Error(t, err)
}
