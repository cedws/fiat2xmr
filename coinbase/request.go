package coinbase

type TxRequest struct {
	Type                        string `json:"type,omitempty"`
	To                          string `json:"to,omitempty"`
	Amount                      string `json:"amount,omitempty"`
	Currency                    string `json:"currency,omitempty"`
	SkipNotifications           bool   `json:"skip_notifications,omitempty"`
	Fee                         string `json:"fee,omitempty"`
	Nonce                       string `json:"nonce,omitempty"`
	ToFinancialInstitution      bool   `json:"to_financial_institution,omitempty"`
	FinancialInstitutionWebsite string `json:"financial_institution_website,omitempty"`
}

type DepositRequest struct {
	Amount        string `json:"amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Commit        bool   `json:"commit,omitempty"`
}

type AdvancedCreateOrderRequest struct {
	ClientOrderID      string `json:"client_order_id,omitempty"`
	ProductID          string `json:"product_id,omitempty"`
	Side               string `json:"side,omitempty"`
	OrderConfiguration struct {
		MarketMarketIOC struct {
			QuoteSize string `json:"quote_size,omitempty"`
			BaseSize  string `json:"base_size,omitempty"`
		} `json:"market_market_ioc,omitempty"`
	} `json:"order_configuration,omitempty"`
}
