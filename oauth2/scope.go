package oauth2

import "strings"

func HierarchicScope(requestScope string, clientScopes []string) bool {
	for _, this := range clientScopes {
		if this == requestScope {
			return true
		}

		if len(this) > len(requestScope) {
			continue
		}

		requestScopes := strings.Split(requestScope, ".")
		haystack := strings.Split(this, ".")
		haystackLen := len(haystack) - 1
		for k, needle := range requestScopes {
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
