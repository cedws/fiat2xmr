package coinbase

import "time"

type AddressesResponse []AddressResponse

type AddressResponse struct {
	ID           string    `json:"id"`
	Address      string    `json:"address"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Network      string    `json:"network"`
	Resource     string    `json:"resource"`
	ResourcePath string    `json:"resource_path"`
}

type ProductResponse struct {
	ProductID                 string  `json:"product_id,omitempty"`
	Price                     float64 `json:"price,omitempty,string"`
	PricePercentageChange24H  string  `json:"price_percentage_change_24h,omitempty"`
	Volume24H                 string  `json:"volume_24h,omitempty"`
	VolumePercentageChange24H string  `json:"volume_percentage_change_24h,omitempty"`
	BaseIncrement             float64 `json:"base_increment,omitempty,string"`
	QuoteIncrement            float64 `json:"quote_increment,omitempty,string"`
	QuoteMinSize              float64 `json:"quote_min_size,omitempty,string"`
	QuoteMaxSize              float64 `json:"quote_max_size,omitempty,string"`
	BaseMinSize               float64 `json:"base_min_size,omitempty,string"`
	BaseMaxSize               float64 `json:"base_max_size,omitempty,string"`
	BaseName                  string  `json:"base_name,omitempty"`
	QuoteName                 string  `json:"quote_name,omitempty"`
	Watched                   bool    `json:"watched,omitempty"`
	IsDisabled                bool    `json:"is_disabled,omitempty"`
	New                       bool    `json:"new,omitempty"`
	Status                    string  `json:"status,omitempty"`
	CancelOnly                bool    `json:"cancel_only,omitempty"`
	LimitOnly                 bool    `json:"limit_only,omitempty"`
	PostOnly                  bool    `json:"post_only,omitempty"`
	TradingDisabled           bool    `json:"trading_disabled,omitempty"`
	AuctionMode               bool    `json:"auction_mode,omitempty"`
	ProductType               string  `json:"product_type,omitempty"`
	QuoteCurrencyID           string  `json:"quote_currency_id,omitempty"`
	BaseCurrencyID            string  `json:"base_currency_id,omitempty"`
}

type AccountsResponse []AccountResponse

type AccountResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Primary  bool   `json:"primary"`
	Type     string `json:"type"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Color        string `json:"color"`
		SortIndex    int    `json:"sort_index"`
		Exponent     int    `json:"exponent"`
		Type         string `json:"type"`
		AddressRegex string `json:"address_regex"`
		AssetID      string `json:"asset_id"`
		Slug         string `json:"slug"`
	} `json:"currency"`
	Balance struct {
		Amount   float64 `json:"amount,string"`
		Currency string  `json:"currency"`
	} `json:"balance"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Resource         string    `json:"resource"`
	ResourcePath     string    `json:"resource_path"`
	AllowDeposits    bool      `json:"allow_deposits"`
	AllowWithdrawals bool      `json:"allow_withdrawals"`
}

type TxResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Amount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"amount"`
	NativeAmount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"native_amount"`
	Description  interface{} `json:"description"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	Resource     string      `json:"resource"`
	ResourcePath string      `json:"resource_path"`
	To           struct {
		Resource string `json:"resource"`
		Email    string `json:"email"`
	} `json:"to"`
	Details struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
	} `json:"details"`
}

type PaymentMethodsResponse []PaymentMethodResponse

type PaymentMethodResponse struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	Currency      string `json:"currency"`
	PrimaryBuy    bool   `json:"primary_buy"`
	PrimarySell   bool   `json:"primary_sell"`
	AllowBuy      bool   `json:"allow_buy"`
	AllowSell     bool   `json:"allow_sell"`
	AllowDeposit  bool   `json:"allow_deposit"`
	AllowWithdraw bool   `json:"allow_withdraw"`
	InstantBuy    bool   `json:"instant_buy"`
	InstantSell   bool   `json:"instant_sell"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Resource      string `json:"resource"`
	ResourcePath  string `json:"resource_path"`
	FiatAccount   struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"fiat_account"`
}

type DepositResponse struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PaymentMethod struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"payment_method"`
	Transaction struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"transaction"`
	Amount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Subtotal struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"subtotal"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Resource     string    `json:"resource"`
	ResourcePath string    `json:"resource_path"`
	Committed    bool      `json:"committed"`
	Fee          struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"fee"`
	PayoutAt string `json:"payout_at"`
}

type AdvancedOrderResponse struct {
	Success         bool   `json:"success"`
	FailureReason   string `json:"failure_reason"`
	OrderID         string `json:"order_id"`
	SuccessResponse struct {
		OrderID       string `json:"order_id"`
		ProductID     string `json:"product_id"`
		Side          string `json:"side"`
		ClientOrderID string `json:"client_order_id"`
	} `json:"success_response"`
	ErrorResponse struct {
		Error                 string `json:"error"`
		Message               string `json:"message"`
		ErrorDetails          string `json:"error_details"`
		PreviewFailureReason  string `json:"preview_failure_reason"`
		NewOrderFailureReason string `json:"new_order_failure_reason"`
	} `json:"error_response"`
	OrderConfiguration struct {
		MarketMarketIoc struct {
			QuoteSize string `json:"quote_size"`
			BaseSize  string `json:"base_size"`
		} `json:"market_market_ioc"`
	} `json:"order_configuration"`
}
