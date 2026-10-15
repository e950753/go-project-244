package formatters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScalar(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"nil", nil, "null"},
		{"string", "text", "text"},
		{"empty string", "", ""},
		{"true", true, "true"},
		{"false", false, "false"},
		{"int", 5, "5"},
		{"int64", int64(7), "7"},
		{"integer float", float64(100500), "100500"},
		{"big float without exponent", float64(1000000), "1000000"},
		{"fractional float", 1.5, "1.5"},
		{"negative float", -0.25, "-0.25"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, scalar(tt.in))
		})
	}
}
