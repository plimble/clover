package scope

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/plimble/clover"
)

//go:generate mockery -name ScopeStorage
type ScopeStorage interface {
	GetScopeByIDs(ids []string) ([]*clover.Scope, error)
}

type ScopeValidator struct {
	storage ScopeStorage
}

func New(storage ScopeStorage) *ScopeValidator {
	return &ScopeValidator{storage}
}

func (s *ScopeValidator) Validate(requestScopes, clientScopes []string) ([]string, error) {
	_, err := s.storage.GetScopeByIDs(requestScopes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(requestScopes) > 0 {
		if len(clientScopes) > 0 {
			if ok := s.check(clientScopes, requestScopes); !ok {
				return nil, errors.New("Invalid scope")
			}
			return requestScopes, nil
		} else {
			return nil, errors.New("Scope is not unsupported")
		}
	} else if len(clientScopes) > 0 {
		return clientScopes, nil
	}

	return nil, errors.New("Scope is not unsupporteds")
}

func (s *ScopeValidator) check(requestScopes, clientScopes []string) bool {
	matched := 0

	for i := 0; i < len(requestScopes); i++ {
		if HierarchicScope(requestScopes[i], clientScopes) {
			matched++
		}
	}

	if matched != len(requestScopes) {
		return false
	}

	return true
}

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
