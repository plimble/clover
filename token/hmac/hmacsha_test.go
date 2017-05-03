package hmac

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFailsWithShortCredentials(t *testing.T) {
	cg := New([]byte("foo"))
	challenge, signature, err := cg.Generate()
	require.NotNil(t, err, "%s", err)
	require.Empty(t, challenge)
	require.Empty(t, signature)
}

func TestGenerate(t *testing.T) {
	cg := New([]byte("12345678901234567890"))

	token, signature, err := cg.Generate()
	require.Nil(t, err, "%s", err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, signature)

	err = cg.Validate(token)
	require.Nil(t, err, "%s", err)

	validateSignature := cg.Signature(token)
	assert.Equal(t, signature, validateSignature)

	// cg.secret = []byte("baz")
	cg = New([]byte("baz"))
	err = cg.Validate(token)
	require.NotNil(t, err, "%s", err)
}

func TestValidateSignatureRejects(t *testing.T) {
	var err error
	cg := New([]byte("12345678901234567890"))
	for k, c := range []string{
		"",
		" ",
		"foo.bar",
		"foo.",
		".foo",
	} {
		err = cg.Validate(c)
		if err != nil {
			assert.NotNil(t, err, "%s", err)
			t.Logf("Passed test case %d", k)
		} else {
			t.Errorf("Failed test case %d", k)
		}
	}
}
