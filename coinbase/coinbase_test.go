package coinbase

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeTime struct{}

func (f fakeTime) Now() time.Time {
	return time.UnixMilli(0)
}

func TestAccessKey(t *testing.T) {
	t.Parallel()

	expectedToken := "123"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedToken, r.Header.Get("CB-ACCESS-KEY"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseURL = srv.URL

	client := NewClient("123", expectedToken)
	_, err := client.SendTransaction("", TxRequest{})
	assert.Nil(t, err)
}

func TestAccessSign(t *testing.T) {
	t.Parallel()

	expectedSig := "bb2ce07c410da073a6b0e89695946ff22ff5f6004016ffbccb28680076911b57"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedSig, r.Header.Get("CB-ACCESS-SIGN"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseURL = srv.URL
	timeNow = fakeTime{}.Now

	client := NewClient("123", "123")
	_, err := client.SendTransaction("", TxRequest{})
	assert.Nil(t, err)
}

func TestAccessTimestamp(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "0", r.Header.Get("CB-ACCESS-TIMESTAMP"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseURL = srv.URL
	timeNow = fakeTime{}.Now

	client := NewClient("123", "123")
	_, err := client.SendTransaction("", TxRequest{})
	assert.Nil(t, err)
}
