package cam

import (
	"strings"
)

func include(data string, lst []string) bool {
	for _, item := range lst {
		if strings.Contains(data, item) {
			return true
		}
	}
	return false
}
