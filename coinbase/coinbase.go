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
	coinbaseV2, _ = url.Parse("https://api.coinbase.com/v2")
	coinbaseV3, _ = url.Parse("https://api.coinbase.com/v3")
	timeNow       = time.Now
)

type Client struct {
	client    *http.Client
	apiKey    string
	apiSecret string
}

func requestV2[T any, U any](c *Client, method, endpoint string, body *T) (*U, error) {
	url := coinbaseV2.JoinPath(endpoint)
	resp, err := request[T, struct {
		Data U
	}](c, method, url, body)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func requestV3[T any, U any](c *Client, method, endpoint string, body *T) (*U, error) {
	url := coinbaseV3.JoinPath(endpoint)
	return request[T, U](c, method, url, body)
}

func request[T any, U any](c *Client, method string, url *url.URL, body *T) (*U, error) {
	bodyReader, bodyWriter := io.Pipe()

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
	// drain body so TCP conn can be reused
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("bad status code %v (%v)", res.StatusCode, http.StatusText(res.StatusCode))
	}
	var decoded U
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return &decoded, nil
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{&http.Client{}, apiKey, apiSecret}
}

func (c *Client) GetAccountByCode(code string) (*AccountResponse, error) {
	path := fmt.Sprintf("/accounts/%v", url.PathEscape(code))

	result, err := requestV2[struct{}, AccountResponse](c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("while getting accounts by code: %w", err)
	}
	return result, nil
}

func (c *Client) GetAccounts() (*AccountsResponse, error) {
	result, err := requestV2[struct{}, AccountsResponse](c, http.MethodGet, "/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("while getting accounts: %w", err)
	}
	return result, nil
}

func (c *Client) GetPaymentMethods() (*PaymentMethodsResponse, error) {
	result, err := requestV2[struct{}, PaymentMethodsResponse](c, http.MethodGet, "/payment-methods", nil)
	if err != nil {
		return nil, fmt.Errorf("while getting payment methods: %w", err)
	}
	return result, nil
}

func (c *Client) GetAddresses(account string) (*AddressesResponse, error) {
	path := fmt.Sprintf("/accounts/%v/addresses", url.PathEscape(account))

	result, err := requestV2[struct{}, AddressesResponse](c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("while getting addresses: %w", err)
	}
	return result, nil
}

func (c *Client) CreateAddress(account string) (*AddressResponse, error) {
	path := fmt.Sprintf("/accounts/%v/addresses", url.PathEscape(account))

	result, err := requestV2[struct{}, AddressResponse](c, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating address: %w", err)
	}
	return result, nil
}

func (c *Client) CreateTransaction(account string, transaction TxRequest) (*TxResponse, error) {
	path := fmt.Sprintf("/accounts/%v/transactions", url.PathEscape(account))

	result, err := requestV2[TxRequest, TxResponse](c, http.MethodPost, path, &transaction)
	if err != nil {
		return nil, fmt.Errorf("while creating transaction: %w", err)
	}
	return result, nil
}

func (c *Client) AdvancedCreateOrder(order AdvancedCreateOrderRequest) (*AdvancedCreateOrderResponse, error) {
	result, err := requestV3[AdvancedCreateOrderRequest, AdvancedCreateOrderResponse](c, http.MethodPost, "/brokerage/orders", &order)
	if err != nil {
		return nil, fmt.Errorf("while creating advanced order: %w", err)
	}
	return result, nil
}
