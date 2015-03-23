package clover

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckScope(t *testing.T) {

	testCases := []struct {
		available []string
		request   []string
		ok        bool
	}{
		{[]string{"1", "2"}, []string{"1", "2"}, true},
		{[]string{"1", "3"}, []string{"1", "2"}, false},
		{[]string{"1", "2", "3"}, []string{"1", "2"}, true},
		{[]string{"1", "2", "3"}, []string{}, true},
		{[]string{}, []string{"1", "2"}, false},
	}

	for _, testCase := range testCases {
		ok := checkScope(testCase.available, testCase.request...)
		assert.Equal(t, testCase.ok, ok)
	}

}
