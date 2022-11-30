package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Overridden in tests.
var (
	coinbaseURL = "https://api.coinbase.com/v2/accounts"
	timeNow     = time.Now
)

type Client struct {
	client  *http.Client
	account string
	token   string
}

type TxRequest struct {
	Type                        string
	To                          string
	Amount                      string
	Currency                    string
	SkipNotifications           bool
	Fee                         string
	Nonce                       string
	ToFinancialInstitution      bool
	FinancialInstitutionWebsite string
}

type TxResponse struct {
	ID     string
	Type   string
	Status string
	Amount struct {
		Amount   string
		Currency string
	}
	NativeAmount struct {
		Amount   string
		Currency string
	}
	Description  string
	CreatedAt    string
	UpdatedAt    string
	Resource     string
	ResourcePath string
	Network      struct {
		Status string
		Hash   string
		Name   string
	}
	To struct {
		Resource string
		Address  string
	}
	Details struct {
		Title    string
		Subtitle string
	}
}

func request[T any, U any](c *Client, endpoint string, body T) (*U, error) {
	bodyReader, bodyWriter := io.Pipe()

	go func() error {
		if err := json.NewEncoder(bodyWriter).Encode(body); err != nil {
			return bodyWriter.CloseWithError(err)
		}
		return bodyWriter.Close()
	}()

	path, err := url.JoinPath(coinbaseURL, c.account, endpoint)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, path, bodyReader)
	if err != nil {
		return nil, err
	}

	timestamp := fmt.Sprintf("%v", timeNow().Unix())
	hmac := hmac.New(sha256.New, []byte(c.token))
	hmac.Write([]byte(timestamp))
	hmac.Write([]byte(http.MethodPost))
	hmac.Write([]byte(endpoint))
	if err := json.NewEncoder(hmac).Encode(body); err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CB-VERSION", "2022-11-28")
	req.Header.Set("CB-ACCESS-KEY", c.token)
	req.Header.Set("CB-ACCESS-SIGN", hex.EncodeToString(hmac.Sum(nil)))
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("bad status code %v (%v)", res.StatusCode, http.StatusText(res.StatusCode))
	}

	var decoded struct {
		Data U
	}
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	return &decoded.Data, nil
}

func NewClient(account, token string) Client {
	return Client{&http.Client{}, account, token}
}

func (c *Client) SendTransaction(transaction TxRequest) (*TxResponse, error) {
	result, err := request[TxRequest, TxResponse](c, "/transactions", transaction)
	if err != nil {
		return nil, fmt.Errorf("coinbase: error sending transaction: %w", err)
	}
	return result, nil
}
