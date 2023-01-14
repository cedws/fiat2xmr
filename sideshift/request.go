package sideshift

type FixedShiftRequest struct {
	SettleAddress string `json:"settleAddress,omitempty"`
	AffiliateID   string `json:"affiliateId,omitempty"`
	QuoteID       string `json:"quoteId,omitempty"`
	RefundAddress string `json:"refundAddress,omitempty"`
}

type QuoteRequest struct {
	DepositCoin    string  `json:"depositCoin,omitempty"`
	DepositNetwork string  `json:"depositNetwork,omitempty"`
	SettleCoin     string  `json:"settleCoin,omitempty"`
	SettleNetwork  string  `json:"settleNetwork,omitempty"`
	DepositAmount  float64 `json:"depositAmount,omitempty,string"`
	SettleAmount   float64 `json:"settleAmount,omitempty,string"`
	AffiliateID    string  `json:"affiliateId,omitempty"`
	CommissionRate float64 `json:"commissionRate,omitempty,string"`
}
