package formatters

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/diff"
)

func TestFormat(t *testing.T) {
	nodes := []diff.Node{
		{Key: "a", Status: diff.Added, NewValue: "v"},
		{Key: "b", Status: diff.Removed, OldValue: float64(1)},
	}

	wantJSON, err := JSON(nodes)
	require.NoError(t, err)

	tests := []struct {
		name, format, want string
	}{
		{"empty format means default", "", Stylish(nodes)},
		{"stylish", "stylish", Stylish(nodes)},
		{"default constant", DefaultFormat, Stylish(nodes)},
		{"plain", "plain", Plain(nodes)},
		{"json", "json", wantJSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(nodes, tt.format)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFormatUnknown(t *testing.T) {
	for _, format := range []string{"xml", "PLAIN", " plain"} { // регистр важен
		t.Run(format, func(t *testing.T) {
			got, err := Format(nil, format)
			require.Error(t, err)
			require.Empty(t, got)
			require.ErrorContains(t, err, "unknown output format")
		})
	}
}
