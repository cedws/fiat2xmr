package sideshift

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Overridden in tests.
var (
	sideshiftURL = "https://sideshift.ai/api/v2"
)

func request[T any, U any](c *Client, method, endpoint string, body T) (*U, error) {
	bodyReader, bodyWriter := io.Pipe()

	go func() error {
		if err := json.NewEncoder(bodyWriter).Encode(body); err != nil {
			return bodyWriter.CloseWithError(err)
		}
		return bodyWriter.Close()
	}()

	path, err := url.JoinPath(sideshiftURL, endpoint)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-sideshift-secret", c.token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		// would be nice to just embed the return type in this struct but alas type parameters cannot be embedded
		var decoded struct {
			Error struct {
				Message string
			}
		}

		if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
			return nil, err
		}

		statusText := http.StatusText(res.StatusCode)
		if len(decoded.Error.Message) != 0 {
			return nil, fmt.Errorf("status %v from server: %v", statusText, decoded.Error.Message)
		}

		return nil, fmt.Errorf("status %v from server", statusText)
	}

	var decoded U
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	return &decoded, nil
}

const (
	StatusWaiting    = "waiting"
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusReview     = "review"
	StatusSettling   = "settling"
	StatusSettled    = "settled"
	StatusRefund     = "refund"
	StatusRefunding  = "refunding"
	StatusRefunded   = "refunded"
	StatusMultiple   = "multiple"
)

type Client struct {
	client *http.Client
	token  string
}

func NewClient(token string) Client {
	return Client{&http.Client{}, token}
}

type WalletAddress string

type VariableShiftRequest struct {
	SettleAddress  WalletAddress
	RefundAddress  WalletAddress
	AffiliateID    string
	DepositCoin    string
	SettleCoin     string
	CommissionRate string
}

type VariableShiftResponse struct {
	ID             string
	CreatedAt      time.Time
	DepositCoin    string
	SettleCoin     string
	DepositNetwork string
	SettleNetwork  string
	DepositAddress string
	SettleAddress  string
	DepositMin     string
	DepositMax     string
	RefundAddress  string
	Type           string
	ExpiresAt      time.Time
	Status         string
}

type GetShiftResponse struct {
	ID                string
	CreatedAt         time.Time
	DepositCoin       string
	SettleCoin        string
	DepositNetwork    string
	SettleNetwork     string
	DepositAddress    string
	SettleAddress     string
	DepositMin        string
	DepositMax        string
	Type              string
	ExpiresAt         time.Time
	Status            string
	UpdatedAt         time.Time
	DepositHash       string
	SettleHash        string
	DepositReceivedAt time.Time
	Rate              string
}

func (c *Client) CreateVariableShift(shift VariableShiftRequest) (*VariableShiftResponse, error) {
	res, err := request[VariableShiftRequest, VariableShiftResponse](c, http.MethodPost, "/shifts/variable", shift)
	if err != nil {
		return nil, fmt.Errorf("sideshift: while creating variable shift: %w", err)
	}

	return res, nil
}

func (c *Client) GetShift(shiftID string) (*GetShiftResponse, error) {
	path, err := url.JoinPath("/shifts", url.PathEscape(shiftID))
	if err != nil {
		panic(err)
	}

	res, err := request[struct{}, GetShiftResponse](c, http.MethodGet, path, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("sideshift: while getting shift: %w", err)
	}

	return res, nil
}
