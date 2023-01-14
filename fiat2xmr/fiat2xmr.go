package fiat2xmr

import (
	"errors"
	"fmt"
	"log"

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
		log.Fatal(err)
	}
	// TODO: somehow check if we have enough to shift before creating an order
	if err := c.assertShiftMinimum(); err != nil {
		log.Fatal(err)
	}

	baseAccount, err := c.cbClient.GetAccountByCode(baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	quote, err := c.ssClient.CreateQuote(sideshift.QuoteRequest{
		DepositCoin:   baseCurrency,
		SettleCoin:    quoteCurrency,
		DepositAmount: baseAccount.Balance.Amount,
	})
	if err != nil {
		log.Fatal(err)
	}

	refundAddress, err := c.getRefundAddress(baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	shift, err := c.ssClient.CreateFixedShift(sideshift.FixedShiftRequest{
		SettleAddress: opts.Address,
		RefundAddress: refundAddress,
		QuoteID:       quote.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.cbClient.CreateTransaction(baseAccount.ID, coinbase.TxRequest{
		Type:     "send",
		To:       shift.DepositAddress,
		Amount:   baseAccount.Balance.Amount,
		Currency: baseCurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	shiftResult, err := c.ssClient.PollShift(shift.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(shiftResult)
}

func (c *Converter) assertShiftMinimum() error {
	baseAccount, err := c.cbClient.GetAccountByCode(baseCurrency)
	if err != nil {
		return err
	}

	pair, err := c.ssClient.GetPair(baseCurrency, quoteCurrency)
	if err != nil {
		return err
	}

	if pair.Min > baseAccount.Balance.Amount {
		log.Fatalf("%v balance too low to initiate shift (minimum %v)", baseCurrency, pair.Min)
	}

	return nil
}

func (c *Converter) createOrder(currency string) error {
	fiatAccount, err := c.cbClient.GetAccountByCode(fiatCurrency)
	if err != nil {
		return err
	}

	if fiatAccount.Balance.Amount > 0 {
		order := coinbase.AdvancedOrderRequest{
			ClientOrderID: uuid.New().String(),
			ProductID:     fmt.Sprintf("%v-%v", currency, fiatCurrency),
			Side:          "BUY",
		}
		order.OrderConfiguration.MarketMarketIOC.QuoteSize = fiatAccount.Balance.Amount

		resp, err := c.cbClient.CreateAdvancedOrder(order)
		if err != nil {
			return err
		}

		if !resp.Success {
			return errors.New(resp.ErrorResponse.Message)
		}
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
