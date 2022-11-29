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

	expectedToken := "1234"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedToken, r.Header.Get("CB-ACCESS-KEY"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseURL = srv.URL

	client, err := NewClient("123", expectedToken)
	assert.Nil(t, err)

	_, err = client.SendTransaction(TxRequest{})
	assert.Nil(t, err)
}

func TestAccessSign(t *testing.T) {
	t.Parallel()

	expectedSig := "dcde1e2bed82538306e8f4a3414690f700368da4f2e853323c425cb12c115e86"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedSig, r.Header.Get("CB-ACCESS-SIGN"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	coinbaseURL = srv.URL
	timeNow = fakeTime{}.Now

	client, err := NewClient("123", "123")
	assert.Nil(t, err)

	_, err = client.SendTransaction(TxRequest{})
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

	client, err := NewClient("123", "123")
	assert.Nil(t, err)

	_, err = client.SendTransaction(TxRequest{})
	assert.Nil(t, err)
}