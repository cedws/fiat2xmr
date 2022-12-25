package coinbase

type BuyOrderRequest struct {
	Amount               string
	Total                string
	Currency             string
	PaymentMethod        string
	AgreeBTCAmountVaries bool
	Commit               bool
	Quote                bool
}

type TxRequest struct {
	Type                        string
	To                          string
	Amount                      string
	Currency                    string
	SkipNotifications           bool
	Fee                         string
	Nonce                       string
	ToFinancialInstitution      bool
	FinancialInstitutionWebsite string
}
