package sideshift

import "time"

type PermissionsResponse struct {
	CreateShift bool `json:"createShift"`
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
	DepositMin     float64 `json:",string"`
	DepositMax     float64 `json:",string"`
	RefundAddress  string
	Type           string
	ExpiresAt      time.Time
	Status         string
}

type ShiftResponse struct {
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
