package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHMACToken(t *testing.T) {
	t.Run("GenerateFailsWithShortCredentials", testGenerateFailsWithShortCredentials)
	t.Run("Generate", testGenerate)
	t.Run("ValidateSignatureRejects", testValidateSignatureRejects)
}

func testGenerateFailsWithShortCredentials(t *testing.T) {
	cg := New([]byte("foo"))
	challenge, err := cg.Generate("", "", "", 0)
	signature := cg.Signature(challenge)
	require.NotNil(t, err, "%s", err)
	require.Empty(t, challenge)
	require.Empty(t, signature)
}

func testGenerate(t *testing.T) {
	cg := New([]byte("12345678901234567890"))

	token, err := cg.Generate("", "", "", 0)
	signature := cg.Signature(token)
	require.Nil(t, err, "%s", err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, signature)

	err = cg.Validate(token)
	require.Nil(t, err, "%s", err)

	validateSignature := cg.Signature(token)
	require.Equal(t, signature, validateSignature)

	// cg.secret = []byte("baz")
	cg = New([]byte("baz"))
	err = cg.Validate(token)
	require.NotNil(t, err, "%s", err)
}

func testValidateSignatureRejects(t *testing.T) {
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
			require.NotNil(t, err, "%s", err)
			t.Logf("Passed test case %d", k)
		} else {
			t.Errorf("Failed test case %d", k)
		}
	}
}
