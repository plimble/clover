package verify

import (
	"net/http"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/plimble/clover/oauth2"
	"github.com/plimble/clover/storage/memory"
	"gopkg.in/gavv/httpexpect.v1"
)

type VerifyHandlerTest struct {
	storage *memory.MemoryStorage
	*httpexpect.Expect
}

func TestVerifyHandler(t *testing.T) {
	ht := &VerifyHandlerTest{
		storage: memory.NewMemoryStorage(),
	}
	h := &VerifyHandler{Logger: zap.L(), Storage: ht.storage}

	config := httpexpect.Config{
		BaseURL: "http://localhost",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(h),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}
	ht.Expect = httpexpect.WithConfig(config)

	t.Run("ValidAccessToken", ht.ValidAccessToken)
}

func (ht *VerifyHandlerTest) ValidAccessToken(t *testing.T) {
	at := &oauth2.AccessToken{
		AccessToken: "a123456",
		ClientID:    "c1",
		UserID:      "u1",
		Expired:     time.Now().Add(time.Second * 1000).Unix(),
		ExpiresIn:   1000,
		Scopes:      []string{"s1", "s2"},
		Extras: map[string]interface{}{
			"e1": "val1",
			"e2": "val1",
		},
	}

	ht.storage.AccessToken["a123456"] = at

	ht.POST("/verify").
		WithHeader("Authorization", "bearer a123456").
		WithFormField("scope", "s1 s2").
		Expect().
		Status(200).
		Body().Equal("")
}
