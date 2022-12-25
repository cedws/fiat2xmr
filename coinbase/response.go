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

type AccountsResponse []AccountResponse

type AccountResponse struct {
	ID       string
	Name     string
	Primary  bool
	Type     string
	Currency struct {
		Code         string
		Name         string
		Color        string
		SortIndex    int
		Exponent     int
		Type         string
		AddressRegex string
		AssetID      string
		Slug         string
	}
	Balance struct {
		Amount   string
		Currency string
	}
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Resource         string
	ResourcePath     string
	AllowDeposits    bool
	AllowWithdrawals bool
	Rewards          struct {
		Apy          string
		FormattedApy string
		Label        string
	}
}

type BuyOrderResponse struct {
	ID            string
	Status        string
	PaymentMethod struct {
		ID           string
		Resource     string
		ResourcePath string
	}
	Transaction struct {
		ID           string
		Resource     string
		ResourcePath string
	}
	Amount struct {
		Amount   string
		Currency string
	}
	Total struct {
		Amount   string
		Currency string
	}
	Subtotal struct {
		Amount   string
		Currency string
	}
	CreatedAt    string
	UpdatedAt    string
	Resource     string
	ResourcePath string
	Committed    bool
	Instant      bool
	Fee          struct {
		Amount   string
		Currency string
	}
	PayoutAt string
}

type TxResponse struct {
	ID     string
	Type   string
	Status string
	Amount struct {
		Amount   string
		Currency string
	}
	NativeAmount struct {
		Amount   string
		Currency string
	}
	Description  string
	CreatedAt    string
	UpdatedAt    string
	Resource     string
	ResourcePath string
	Network      struct {
		Status string
		Hash   string
		Name   string
	}
	To struct {
		Resource string
		Address  string
	}
	Details struct {
		Title    string
		Subtitle string
	}
}

type PaymentMethodsResponse []PaymentMethodResponse

type PaymentMethodResponse struct {
	ID            string
	Type          string
	Name          string
	Currency      string
	PrimaryBuy    bool
	PrimarySell   bool
	AllowBuy      bool
	AllowSell     bool
	AllowDeposit  bool
	AllowWithdraw bool
	InstantBuy    bool
	InstantSell   bool
	CreatedAt     string
	UpdatedAt     string
	Resource      string
	ResourcePath  string
	FiatAccount   struct {
		ID           string
		Resource     string
		ResourcePath string
	}
}
