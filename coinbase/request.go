package coinbase

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
