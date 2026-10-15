package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/parsers"
)

func TestParseJSON(t *testing.T) {
	got, err := parsers.Parse([]byte(`{"a": 50, "b": {"c": "x"}, "d": null, "e": [1, 2]}`), ".json")
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"a": float64(50), // числа из JSON — float64
		"b": map[string]any{"c": "x"},
		"d": nil,
		"e": []any{float64(1), float64(2)},
	}, got)
}

func TestParseYAML(t *testing.T) {
	data := []byte("a: 50\nb:\n  c: x\nd: null\n")
	want := map[string]any{
		"a": float64(50),
		"b": map[string]any{"c": "x"},
		"d": nil,
	}

	for _, ext := range []string{".yml", ".yaml"} {
		t.Run(ext, func(t *testing.T) {
			got, err := parsers.Parse(data, ext)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
}

func TestParseExtensionIsCaseInsensitive(t *testing.T) {
	for _, ext := range []string{".JSON", ".Json"} {
		got, err := parsers.Parse([]byte(`{"a": 1}`), ext)
		require.NoError(t, err)
		require.Equal(t, map[string]any{"a": float64(1)}, got)
	}

	got, err := parsers.Parse([]byte("a: 1"), ".YML")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"a": float64(1)}, got)
}

func TestParseEmptyDocuments(t *testing.T) {
	// null в JSON и пустой YAML дают пустую мапу, а не nil
	got, err := parsers.Parse([]byte("null"), ".json")
	require.NoError(t, err)
	require.Equal(t, map[string]any{}, got)

	got, err = parsers.Parse([]byte(""), ".yml")
	require.NoError(t, err)
	require.Equal(t, map[string]any{}, got)

	got, err = parsers.Parse([]byte("{}"), ".json")
	require.NoError(t, err)
	require.Equal(t, map[string]any{}, got)
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		data string
		ext  string
	}{
		{"unsupported extension", `{}`, ".txt"},
		{"no extension", `{}`, ""},
		{"broken json", `{"a": `, ".json"},
		{"empty json file", ``, ".json"},
		{"json array on top level", `[1, 2]`, ".json"},
		{"broken yaml", "a: [1, 2", ".yml"},
		{"yaml list on top level", "- a\n- b\n", ".yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsers.Parse([]byte(tt.data), tt.ext)
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func TestParseYAMLNormalizesNested(t *testing.T) {
	got, err := parsers.Parse([]byte("a:\n  b: 1\nc: [1, 2]\n"), ".yml")
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"a": map[string]any{"b": float64(1)},
		"c": []any{float64(1), float64(2)},
	}, got)
}
