package sideshift

import "time"

type PermissionsResponse struct {
	CreateShift bool `json:"createShift"`
}

type PairResponse struct {
	Min            float64 `json:"min,omitempty,string"`
	Max            float64 `json:"max,omitempty,string"`
	Rate           float64 `json:"rate,omitempty,string"`
	DepositCoin    string  `json:"depositCoin,omitempty"`
	SettleCoin     string  `json:"settleCoin,omitempty"`
	DepositNetwork string  `json:"depositNetwork,omitempty"`
	SettleNetwork  string  `json:"settleNetwork,omitempty"`
}

type FixedShiftResponse struct {
	ID             string    `json:"id,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	DepositCoin    string    `json:"depositCoin,omitempty"`
	SettleCoin     string    `json:"settleCoin,omitempty"`
	DepositNetwork string    `json:"depositNetwork,omitempty"`
	SettleNetwork  string    `json:"settleNetwork,omitempty"`
	DepositAddress string    `json:"depositAddress,omitempty"`
	SettleAddress  string    `json:"settleAddress,omitempty"`
	DepositMin     string    `json:"depositMin,omitempty"`
	DepositMax     string    `json:"depositMax,omitempty"`
	RefundAddress  string    `json:"refundAddress,omitempty"`
	Type           string    `json:"type,omitempty"`
	QuoteID        string    `json:"quoteId,omitempty"`
	DepositAmount  string    `json:"depositAmount,omitempty"`
	SettleAmount   string    `json:"settleAmount,omitempty"`
	ExpiresAt      time.Time `json:"expiresAt,omitempty"`
	Status         string    `json:"status,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
	Rate           string    `json:"rate,omitempty"`
}

type ShiftResponse struct {
	ID                string    `json:"id,omitempty"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	DepositCoin       string    `json:"depositCoin,omitempty"`
	SettleCoin        string    `json:"settleCoin,omitempty"`
	DepositNetwork    string    `json:"depositNetwork,omitempty"`
	SettleNetwork     string    `json:"settleNetwork,omitempty"`
	DepositAddress    string    `json:"depositAddress,omitempty"`
	SettleAddress     string    `json:"settleAddress,omitempty"`
	DepositMin        float64   `json:"depositMin,omitempty,string"`
	DepositMax        float64   `json:"depositMax,omitempty,string"`
	Type              string    `json:"type,omitempty"`
	ExpiresAt         time.Time `json:"expiresAt,omitempty"`
	Status            string    `json:"status,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
	DepositHash       string    `json:"depositHash,omitempty"`
	SettleHash        string    `json:"settleHash,omitempty"`
	DepositReceivedAt time.Time `json:"depositReceivedAt,omitempty"`
	Rate              float64   `json:"rate,omitempty,string"`
}

type QuoteResponse struct {
	ID             string    `json:"id,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	DepositCoin    string    `json:"depositCoin,omitempty"`
	SettleCoin     string    `json:"settleCoin,omitempty"`
	DepositNetwork string    `json:"depositNetwork,omitempty"`
	SettleNetwork  string    `json:"settleNetwork,omitempty"`
	ExpiresAt      time.Time `json:"expiresAt,omitempty"`
	DepositAmount  float64   `json:"depositAmount,omitempty,string"`
	SettleAmount   float64   `json:"settleAmount,omitempty,string"`
	Rate           float64   `json:"rate,omitempty,string"`
	AffiliateID    string    `json:"affiliateId,omitempty"`
}
