package sideshift

type VariableShiftRequest struct {
	SettleAddress  string `json:"settleAddress,omitempty"`
	RefundAddress  string `json:"refundAddress,omitempty"`
	AffiliateID    string `json:"affiliateId,omitempty"`
	DepositCoin    string `json:"depositCoin,omitempty"`
	SettleCoin     string `json:"settleCoin,omitempty"`
	CommissionRate string `json:"commissionRate,omitempty"`
}
