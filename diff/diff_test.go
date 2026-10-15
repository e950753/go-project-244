package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/diff"
)

type obj = map[string]any

func TestBuildEmpty(t *testing.T) {
	require.Empty(t, diff.Build(nil, nil))
	require.Empty(t, diff.Build(obj{}, obj{}))
}

func TestBuild(t *testing.T) {
	tests := []struct {
		name string
		a, b obj
		want []diff.Node
	}{
		{
			name: "added",
			a:    obj{}, b: obj{"k": 1},
			want: []diff.Node{{Key: "k", Status: diff.Added, NewValue: 1}},
		},
		{
			name: "removed",
			a:    obj{"k": 1}, b: obj{},
			want: []diff.Node{{Key: "k", Status: diff.Removed, OldValue: 1}},
		},
		{
			name: "unchanged",
			a:    obj{"k": "v"}, b: obj{"k": "v"},
			want: []diff.Node{{Key: "k", Status: diff.Unchanged, OldValue: "v", NewValue: "v"}},
		},
		{
			name: "updated",
			a:    obj{"k": 1}, b: obj{"k": 2},
			want: []diff.Node{{Key: "k", Status: diff.Updated, OldValue: 1, NewValue: 2}},
		},
		{
			name: "nil to value is updated",
			a:    obj{"k": nil}, b: obj{"k": false},
			want: []diff.Node{{Key: "k", Status: diff.Updated, OldValue: nil, NewValue: false}},
		},
		{
			name: "nil in both is unchanged",
			a:    obj{"k": nil}, b: obj{"k": nil},
			want: []diff.Node{{Key: "k", Status: diff.Unchanged}},
		},
		{
			name: "removed key with nil value",
			a:    obj{"k": nil}, b: obj{},
			want: []diff.Node{{Key: "k", Status: diff.Removed}},
		},
		{
			name: "object to scalar is updated, not nested",
			a:    obj{"k": obj{"x": 1}}, b: obj{"k": "str"},
			want: []diff.Node{{Key: "k", Status: diff.Updated, OldValue: obj{"x": 1}, NewValue: "str"}},
		},
		{
			name: "scalar to object is updated, not nested",
			a:    obj{"k": "str"}, b: obj{"k": obj{"x": 1}},
			want: []diff.Node{{Key: "k", Status: diff.Updated, OldValue: "str", NewValue: obj{"x": 1}}},
		},
		{
			name: "equal slices",
			a:    obj{"k": []any{1, 2}}, b: obj{"k": []any{1, 2}},
			want: []diff.Node{{Key: "k", Status: diff.Unchanged, OldValue: []any{1, 2}, NewValue: []any{1, 2}}},
		},
		{
			name: "different slices",
			a:    obj{"k": []any{1, 2}}, b: obj{"k": []any{2, 1}},
			want: []diff.Node{{Key: "k", Status: diff.Updated, OldValue: []any{1, 2}, NewValue: []any{2, 1}}},
		},
		{
			name: "empty nested objects",
			a:    obj{"k": obj{}}, b: obj{"k": obj{}},
			want: []diff.Node{{Key: "k", Status: diff.Nested, Children: []diff.Node{}}},
		},
		{
			name: "nested recursion",
			a:    obj{"n": obj{"a": 1, "b": 2}}, b: obj{"n": obj{"a": 1, "c": 3}},
			want: []diff.Node{{
				Key: "n", Status: diff.Nested,
				Children: []diff.Node{
					{Key: "a", Status: diff.Unchanged, OldValue: 1, NewValue: 1},
					{Key: "b", Status: diff.Removed, OldValue: 2},
					{Key: "c", Status: diff.Added, NewValue: 3},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, diff.Build(tt.a, tt.b))
		})
	}
}

func TestBuildNilMaps(t *testing.T) {
	got := diff.Build(nil, obj{"k": 1})
	require.Equal(t, []diff.Node{{Key: "k", Status: diff.Added, NewValue: 1}}, got)

	got = diff.Build(obj{"k": 1}, nil)
	require.Equal(t, []diff.Node{{Key: "k", Status: diff.Removed, OldValue: 1}}, got)
}

func TestBuildSortsKeys(t *testing.T) {
	got := diff.Build(obj{"b": 1, "c": 1}, obj{"a": 1, "c": 1})

	keys := make([]string, 0, len(got))
	for _, n := range got {
		keys = append(keys, n.Key)
	}
	require.Equal(t, []string{"a", "b", "c"}, keys)
}
