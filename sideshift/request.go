package sideshift

type FixedShiftRequest struct {
	SettleAddress string `json:"settleAddress,omitempty"`
	QuoteID       string `json:"quoteId,omitempty"`
	RefundAddress string `json:"refundAddress,omitempty"`
}

type QuoteRequest struct {
	DepositCoin   string  `json:"depositCoin,omitempty"`
	SettleCoin    string  `json:"settleCoin,omitempty"`
	DepositAmount float64 `json:"depositAmount,omitempty,string"`
	SettleAmount  float64 `json:"settleAmount,omitempty,string"`
}
