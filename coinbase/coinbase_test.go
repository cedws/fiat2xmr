package coinbase

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeTime struct{}

func (f fakeTime) Now() time.Time {
	return time.UnixMilli(0)
}

func TestAccessKey(t *testing.T) {
	expectedToken := "123"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedToken, r.Header.Get("CB-ACCESS-KEY"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseV2, _ = url.Parse(srv.URL)

	client := NewClient("123", expectedToken)
	_, err := client.CreateTransaction("", TxRequest{})
	assert.Nil(t, err)
}

func TestAccessSign(t *testing.T) {
	expectedSig := "3e8d467949bb98bc2c9f620aebc822930e2ffbcfd77d09b8b78461a18f5acdf2"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedSig, r.Header.Get("CB-ACCESS-SIGN"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseV2, _ = url.Parse(srv.URL)
	timeNow = fakeTime{}.Now

	client := NewClient("123", "123")
	_, err := client.CreateTransaction("", TxRequest{})
	assert.Nil(t, err)
}

func TestAccessTimestamp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "0", r.Header.Get("CB-ACCESS-TIMESTAMP"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseV2, _ = url.Parse(srv.URL)
	timeNow = fakeTime{}.Now

	client := NewClient("123", "123")
	_, err := client.CreateTransaction("", TxRequest{})
	assert.Nil(t, err)
}
