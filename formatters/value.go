package formatters

import (
	"fmt"
	"strconv"
)

func scalar(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		return x
	case float64:
		// числа из JSON; 'f' не даёт 1e+06
		return strconv.FormatFloat(x, 'f', -1, 64)
	default:
		return fmt.Sprint(x)
	}
}
