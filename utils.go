package supersql

import (
	"fmt"
	"strings"
)

func coerceToString(input interface{}) (t string) {
	switch val := input.(type) {
	case *SqlTable:
		t = val.name
	case string:
		t = val
	}
	return
}

func countAndReplacePlaceholders(ssql string) string {
	if placeholders := strings.Count(ssql, "?"); placeholders > 0 {
		for i := 0; i < placeholders; i++ {
			//remember $params start at $1 not $0 so offset index here
			ssql = strings.Replace(ssql, "?", fmt.Sprintf("$%d", i+1), 1)
		}
	}
	return ssql
}
