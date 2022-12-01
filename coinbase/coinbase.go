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
	coinbaseURL = "https://api.coinbase.com/v2"
	timeNow     = time.Now
)

type Client struct {
	client    *http.Client
	apiKey    string
	apiSecret string
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

type Accounts []Account

type Account struct {
	ID       string
	Name     string
	Primary  bool
	Type     string
	Currency struct {
		Code         string
		Name         string
		Color        string
		SortIndex    int
		Exponent     int
		Type         string
		AddressRegex string
		AssetID      string
		Slug         string
	}
	Balance struct {
		Amount   string
		Currency string
	}
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Resource         string
	ResourcePath     string
	AllowDeposits    bool
	AllowWithdrawals bool
	Rewards          struct {
		Apy          string
		FormattedApy string
		Label        string
	}
}

func request[T any, U any](c *Client, method, endpoint string, body *T) (*U, error) {
	bodyReader, bodyWriter := io.Pipe()

	coinbaseURL, err := url.Parse(coinbaseURL)
	url := coinbaseURL.JoinPath(endpoint)
	if err != nil {
		// should not happen
		panic(err)
	}

	timestamp := fmt.Sprintf("%v", timeNow().Unix())
	hmac := hmac.New(sha256.New, []byte(c.apiSecret))
	hmac.Write([]byte(timestamp))
	hmac.Write([]byte(method))
	hmac.Write([]byte(url.Path))

	if body != nil && method != http.MethodGet {
		go func() error {
			if err := json.NewEncoder(bodyWriter).Encode(body); err != nil {
				return bodyWriter.CloseWithError(err)
			}

			return bodyWriter.Close()
		}()

		if err := json.NewEncoder(hmac).Encode(body); err != nil {
			return nil, err
		}
	} else {
		bodyWriter.Close()
	}

	req, err := http.NewRequest(method, url.String(), bodyReader)
	if err != nil {
		// should not happen
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CB-VERSION", "2022-11-28")
	req.Header.Set("CB-ACCESS-KEY", c.apiKey)
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

func NewClient(apiKey, apiSecret string) Client {
	return Client{&http.Client{}, apiKey, apiSecret}
}

func (c *Client) GetAccountByCode(code string) (*Account, error) {
	path := fmt.Sprintf("/accounts/%v", url.PathEscape(code))

	result, err := request[struct{}, Account](c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("while getting accounts: %w", err)
	}
	return result, nil
}

func (c *Client) GetAccounts() (*Accounts, error) {
	result, err := request[struct{}, Accounts](c, http.MethodGet, "/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("while getting accounts: %w", err)
	}
	return result, nil
}

func (c *Client) SendTransaction(account string, transaction TxRequest) (*TxResponse, error) {
	path := fmt.Sprintf("/accounts/%v/transactions", url.PathEscape(account))

	result, err := request[TxRequest, TxResponse](c, http.MethodPost, path, &transaction)
	if err != nil {
		return nil, fmt.Errorf("while sending transaction: %w", err)
	}
	return result, nil
}
