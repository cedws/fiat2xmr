package fiat2xmr

import (
	"fmt"
	"math"

	"github.com/apex/log"
	"github.com/cedws/fiat2xmr/coinbase"
	"github.com/cedws/fiat2xmr/sideshift"
	"github.com/google/uuid"
)

const (
	fiatCurrency  = "GBP"
	baseCurrency  = "LTC"
	quoteCurrency = "XMR"
)

type Opts struct {
	CoinbaseKey     string
	CoinbaseSecret  string
	SideShiftSecret string
	Address         string
}

type Converter struct {
	ssClient *sideshift.Client
	cbClient *coinbase.Client
}

func NewConverter(ssClient *sideshift.Client, cbClient *coinbase.Client) *Converter {
	return &Converter{ssClient, cbClient}
}

func (c *Converter) Convert(opts Opts) {
	if err := c.createOrder(baseCurrency); err != nil {
		log.Fatalf("%v", err)
	}

	baseAccount, err := c.cbClient.GetAccountByCode(baseCurrency)
	if err != nil {
		log.Fatalf("%v", err)
	}

	quote, err := c.ssClient.CreateQuote(sideshift.QuoteRequest{
		DepositCoin:   baseCurrency,
		SettleCoin:    quoteCurrency,
		DepositAmount: baseAccount.Balance.Amount,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Infof("shift quote price is %v, expires at %v", quote.Rate, quote.ExpiresAt)

	refundAddress, err := c.getRefundAddress(baseCurrency)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Infof("using %v as base refund address", refundAddress)

	log.Infof("creating fixed shift")
	shift, err := c.ssClient.CreateFixedShift(sideshift.FixedShiftRequest{
		SettleAddress: opts.Address,
		RefundAddress: refundAddress,
		QuoteID:       quote.ID,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Infof("sending %v %v to shift address %v", baseAccount.Balance.Amount, baseCurrency, shift.DepositAddress)
	_, err = c.cbClient.CreateTransaction(baseAccount.ID, coinbase.TxRequest{
		Type:     "send",
		To:       shift.DepositAddress,
		Amount:   baseAccount.Balance.Amount,
		Currency: baseCurrency,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Info("waiting for shift completion")
	shiftResult, err := c.ssClient.PollShift(shift.ID)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Infof("%+v", shiftResult)
}

func (c *Converter) getBalance(currency string) (float64, error) {
	account, err := c.cbClient.GetAccountByCode(currency)
	if err != nil {
		return 0, err
	}

	return account.Balance.Amount, nil
}

func (c *Converter) createOrder(currency string) error {
	productID := fmt.Sprintf("%v-%v", currency, fiatCurrency)
	log.Infof("using product %v", productID)

	product, err := c.cbClient.GetProduct(productID)
	if err != nil {
		return err
	}
	if product.TradingDisabled {
		return fmt.Errorf("trading for product %v is disabled", productID)
	}

	pair, err := c.ssClient.GetPair(baseCurrency, quoteCurrency)
	if err != nil {
		return err
	}
	log.Infof("shift minimum is %v, maximum is %v", pair.Min, pair.Max)

	fiatBalance, err := c.getBalance(fiatCurrency)
	if err != nil {
		return err
	}
	log.Infof("fiat balance is %v", fiatBalance)

	// quote means fiat here thanks to coinbase inverting things
	if fiatBalance > 0 && fiatBalance > product.QuoteMinSize {
		// clamp amount to maximum order size for the millionaires
		orderVolumeFiat := math.Min(fiatBalance, product.QuoteMaxSize)

		baseBalance, err := c.getBalance(baseCurrency)
		if err != nil {
			return err
		}
		log.Infof("base balance is %v", baseBalance)
		// estimate if we'll have enough to shift if we place a market order
		if baseBalance+(product.Price/orderVolumeFiat) < pair.Min {
			return fmt.Errorf("%v balance too low to initiate shift (minimum %v)", baseCurrency, pair.Min)
		}

		log.Infof("placing order for %v %v", orderVolumeFiat, baseCurrency)
		order := coinbase.AdvancedOrderRequest{
			ClientOrderID: uuid.New().String(),
			ProductID:     productID,
			Side:          "BUY",
		}
		order.OrderConfiguration.MarketMarketIOC.QuoteSize = orderVolumeFiat

		resp, err := c.cbClient.CreateAdvancedOrder(order)
		if err != nil {
			return err
		}
		if !resp.Success {
			return fmt.Errorf("advanced order failed: %v", resp.ErrorResponse.Message)
		}

		log.Info("order succeeded")
	}

	baseBalance, err := c.getBalance(baseCurrency)
	if err != nil {
		return err
	}
	log.Infof("base balance is %v", baseBalance)
	// additional check before we start the shift just in case the price moved since the pre-flight check
	if pair.Min > baseBalance {
		return fmt.Errorf("%v balance too low to initiate shift (minimum %v)", baseCurrency, pair.Min)
	}
	if pair.Max < baseBalance {
		return fmt.Errorf("%v balance too high to initiate shift (maximum %v)", baseCurrency, pair.Max)
	}

	return nil
}

func (c *Converter) getRefundAddress(currency string) (string, error) {
	addresses, err := c.cbClient.GetAddresses(baseCurrency)
	if err != nil {
		return "", err
	}

	var refundAddress string
	if len(*addresses) > 0 {
		refundAddress = (*addresses)[0].Address
	} else {
		address, err := c.cbClient.CreateAddress(baseCurrency)
		if err != nil {
			return "", err
		}
		refundAddress = address.Address
	}

	return refundAddress, nil
}
