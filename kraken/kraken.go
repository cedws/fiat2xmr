package kraken

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func request[T encoding.TextMarshaler, U any](c *Client, params T, path string) (*U, error) {
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	body, _ := params.MarshalText()
	body = append(body, []byte("&nonce="+nonce)...)

	req, err := http.NewRequest(http.MethodPost, "https://api.kraken.com"+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	sha := sha256.New()
	sha.Write([]byte(nonce))
	sha.Write(body)
	shasum := sha.Sum(nil)

	mac := hmac.New(sha512.New, c.privateKey)
	mac.Write([]byte(path))
	mac.Write(shasum)
	macsum := mac.Sum(nil)

	req.Header.Set("API-Key", c.apiKey)
	req.Header.Set("API-Sign", base64.StdEncoding.EncodeToString(macsum))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decodedBody U
	if err := json.NewDecoder(resp.Body).Decode(&decodedBody); err != nil {
		return nil, err
	}

	return &decodedBody, nil
}

type Client struct {
	apiKey     string
	privateKey []byte
	httpClient http.Client
}

func NewClient(apiKey, privateKey string) (*Client, error) {
	privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return &Client{apiKey, privateKeyDecoded, http.Client{}}, nil
}

type AddOrderRequest struct {
	Pair      string
	Type      string
	OrderType string
	Volume    string
	Oflags    []string
	Validate  bool
}

type AddOrderResponse struct {
	Error  []string
	Result struct {
		Descr struct {
			Order string
			Close string
		}
		Txid []string
	}
}

func (a AddOrderRequest) MarshalText() ([]byte, error) {
	params := url.Values{}
	params.Set("pair", a.Pair)
	params.Set("type", a.Type)
	params.Set("ordertype", a.OrderType)
	params.Set("volume", a.Volume)
	params.Set("oflags", strings.Join(a.Oflags, ","))
	params.Set("validate", strconv.FormatBool(a.Validate))

	return []byte(params.Encode()), nil
}

func (c *Client) AddOrder() (*AddOrderResponse, error) {
	order := AddOrderRequest{
		Pair:      "LTCGBP",
		Type:      "buy",
		OrderType: "market",
		Volume:    "1.0",
		Oflags:    []string{"fciq"},
		Validate:  true,
	}

	return request[AddOrderRequest, AddOrderResponse](c, order, "/0/private/AddOrder")
}
