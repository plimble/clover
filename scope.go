package clover

import (
	"strings"
)

func CheckScope(scopes, needles []string) bool {
	matched := 0

	for i := 0; i < len(needles); i++ {
		if hierarchicScope(scopes, needles[i]) {
			matched++
		}
	}

	if matched != len(scopes) {
		return false
	}

	return true
}

func hierarchicScope(haystack []string, needle string) bool {
	for _, this := range haystack {
		if this == needle {
			return true
		}

		if len(this) > len(needle) {
			continue
		}

		needles := strings.Split(needle, ".")
		haystack := strings.Split(this, ".")
		haystackLen := len(haystack) - 1
		for k, needle := range needles {
			if haystackLen < k {
				return true
			}

			current := haystack[k]
			if current != needle {
				break
			}
		}
	}

	return false
}
