package code

import (
	"code/diff"
	"code/formatters"
	"code/parsers"
	"fmt"
	"os"
	"path/filepath"
)

func GenDiff(path1, path2, format string) (string, error) {
	data1, err := readFile(path1)
	if err != nil {
		return "", err
	}
	data2, err := readFile(path2)
	if err != nil {
		return "", err
	}

	return formatters.Format(diff.Build(data1, data2), format)
}

func readFile(path string) (map[string]any, error) {
	content, err := os.ReadFile(path) // работает и с относительными, и с абсолютными путями
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	data, err := parsers.Parse(content, filepath.Ext(path))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return data, nil
}
