package code_test

import (
	"code"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func fixture(name string) string {
	return filepath.Join("testdata", "fixture", name)
}

func expected(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(fixture(name))
	require.NoError(t, err)
	return strings.TrimRight(string(data), "\r\n")
}

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name, file1, file2, format, result string
	}{
		{"json stylish", "nested1.json", "nested2.json", "stylish", "result_stylish.txt"},
		{"yaml stylish", "nested1.yml", "nested2.yml", "stylish", "result_stylish.txt"},
		{"default format", "nested1.json", "nested2.json", "", "result_stylish.txt"},
		{"json plain", "nested1.json", "nested2.json", "plain", "result_plain.txt"},
		{"yaml plain", "nested1.yml", "nested2.yml", "plain", "result_plain.txt"},
		{"mixed json+yaml", "nested1.json", "nested2.yml", "plain", "result_plain.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := code.GenDiff(fixture(tt.file1), fixture(tt.file2), tt.format)
			require.NoError(t, err)
			require.Equal(t, expected(t, tt.result), got)
		})
	}
}

func TestGenDiffAbsolutePath(t *testing.T) {
	abs, err := filepath.Abs(fixture("nested1.json"))
	require.NoError(t, err)
	_, err = code.GenDiff(abs, abs, "stylish")
	require.NoError(t, err)
}

func TestGenDiffErrors(t *testing.T) {
	_, err := code.GenDiff("nope.json", fixture("nested2.json"), "stylish")
	require.Error(t, err)

	_, err = code.GenDiff(fixture("result_plain.txt"), fixture("nested2.json"), "stylish")
	require.Error(t, err)

	_, err = code.GenDiff(fixture("nested1.json"), fixture("nested2.json"), "xml")
	require.Error(t, err)
}

func TestGenDiffJSON(t *testing.T) {
	for _, pair := range [][2]string{
		{"nested1.json", "nested2.json"},
		{"nested1.yml", "nested2.yml"},
		{"nested1.json", "nested2.yml"},
	} {
		t.Run(pair[0]+"+"+pair[1], func(t *testing.T) {
			got, err := code.GenDiff(fixture(pair[0]), fixture(pair[1]), "json")
			require.NoError(t, err)
			require.True(t, json.Valid([]byte(got)))

			var tree []map[string]any
			require.NoError(t, json.Unmarshal([]byte(got), &tree))

			// верхний уровень отсортирован: common, group1, group2, group3
			keys := make([]string, 0, len(tree))
			for _, n := range tree {
				keys = append(keys, n["key"].(string))
			}
			require.Equal(t, []string{"common", "group1", "group2", "group3"}, keys)

			require.Equal(t, "nested", tree[0]["status"])
			require.Equal(t, "removed", tree[2]["status"])
			require.Equal(t, "added", tree[3]["status"])
		})
	}
}

func TestGenDiffJSONKeepsNull(t *testing.T) {
	got, err := code.GenDiff(fixture("nested1.json"), fixture("nested2.json"), "json")
	require.NoError(t, err)

	var tree []map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &tree))

	var setting3 map[string]any
	for _, c := range tree[0]["children"].([]any) {
		if m := c.(map[string]any); m["key"] == "setting3" {
			setting3 = m
		}
	}
	require.NotNil(t, setting3)
	require.Equal(t, "updated", setting3["status"])
	require.Equal(t, true, setting3["oldValue"])
	require.Contains(t, setting3, "newValue") // ключ есть, значение null
	require.Nil(t, setting3["newValue"])
}
