package coinbase

type TxRequest struct {
	Type                        string `json:"type"`
	To                          string `json:"to"`
	Amount                      string `json:"amount"`
	Currency                    string `json:"currency"`
	SkipNotifications           bool   `json:"skip_notifications"`
	Fee                         string `json:"fee"`
	Nonce                       string `json:"nonce"`
	ToFinancialInstitution      bool   `json:"to_financial_institution"`
	FinancialInstitutionWebsite string `json:"financial_institution_website"`
}

type AdvancedCreateOrderRequest struct {
	ClientOrderID      string `json:"client_order_id"`
	ProductID          string `json:"product_id"`
	Side               string `json:"side"`
	OrderConfiguration struct {
		MarketMarketIOC struct {
			QuoteSize string `json:"quote_size"`
			BaseSize  string `json:"base_size"`
		} `json:"market_market_ioc"`
	} `json:"order_configuration"`
}
