package sideshift

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

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

// Overridden in tests.
var (
	sideshiftV2 = "https://sideshift.ai/api/v2"
)

func request[T any, U any](c *Client, method, endpoint string, body *T) (*U, error) {
	bodyReader, bodyWriter := io.Pipe()

	if body != nil && method != http.MethodGet {
		go func() error {
			if err := json.NewEncoder(bodyWriter).Encode(body); err != nil {
				return bodyWriter.CloseWithError(err)
			}
			return bodyWriter.Close()
		}()
	} else {
		bodyWriter.Close()
	}

	path, err := url.JoinPath(sideshiftV2, endpoint)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-sideshift-secret", c.apiSecret)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// drain body so TCP conn can be reused
	defer io.Copy(io.Discard, res.Body)

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

		if len(decoded.Error.Message) != 0 {
			return nil, fmt.Errorf("bad status code %v: %v", res.StatusCode, decoded.Error.Message)
		}

		return nil, fmt.Errorf("bad status code %v", res.StatusCode)
	}

	var decoded U
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	return &decoded, nil
}

type Client struct {
	client    *http.Client
	apiSecret string
}

func NewClient(apiSecret string) Client {
	return Client{&http.Client{}, apiSecret}
}

func (c *Client) CanShift() (bool, error) {
	perms, err := c.GetPermissions()
	if err != nil {
		return false, err
	}

	return perms.CreateShift, nil
}

func (c *Client) GetPermissions() (*PermissionsResponse, error) {
	res, err := request[struct{}, PermissionsResponse](c, http.MethodGet, "/permissions", nil)
	if err != nil {
		return nil, fmt.Errorf("while getting permissions: %w", err)
	}

	return res, nil
}

func (c *Client) CreateVariableShift(shift VariableShiftRequest) (*VariableShiftResponse, error) {
	res, err := request[VariableShiftRequest, VariableShiftResponse](c, http.MethodPost, "/shifts/variable", &shift)
	if err != nil {
		return nil, fmt.Errorf("while creating variable shift: %w", err)
	}

	return res, nil
}

func (c *Client) PollShift(shiftID string) (shift *ShiftResponse, err error) {
Loop:
	for range time.Tick(10 * time.Second) {
		shift, err = c.GetShift(shiftID)
		if err != nil {
			return nil, err
		}

		if time.Now().After(shift.ExpiresAt) {
			// TODO: return error?
			return shift, nil
		}

		switch status := shift.Status; status {
		case StatusWaiting, StatusPending, StatusProcessing, StatusSettling:
			continue
		case StatusSettled:
			// OK, all done
			break Loop
		default:
			return nil, fmt.Errorf("unknown shift status %v", status)
		}
	}

	return shift, nil
}

func (c *Client) GetShift(shiftID string) (*ShiftResponse, error) {
	path, err := url.JoinPath("/shifts", url.PathEscape(shiftID))
	if err != nil {
		panic(err)
	}

	res, err := request[struct{}, ShiftResponse](c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("while getting shift: %w", err)
	}

	return res, nil
}
