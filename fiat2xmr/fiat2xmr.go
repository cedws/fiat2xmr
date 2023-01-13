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
	baseCurrency  = "LTC"
	quoteCurrency = "XMR"
)

type Opts struct {
	CoinbaseKey     string
	CoinbaseSecret  string
	SideShiftSecret string
	Address         string
	VolumeBase      float64
}

func Convert(opts Opts) {
	ssClient := sideshift.NewClient(opts.SideShiftSecret)
	if canShift, err := ssClient.CanShift(); !canShift || err != nil {
		log.Fatal("sideshift account is unable to create shifts")
	}

	cbClient := coinbase.NewClient(opts.CoinbaseKey, opts.CoinbaseSecret)
	if err := createOrder(cbClient, baseCurrency); err != nil {
		log.Fatal(err)
	}

	refundAddress, err := getRefundAddress(cbClient, baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	shift, err := ssClient.CreateVariableShift(sideshift.VariableShiftRequest{
		SettleAddress: opts.Address,
		RefundAddress: refundAddress,
		DepositCoin:   baseCurrency,
		SettleCoin:    quoteCurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	baseAccount, err := cbClient.GetAccountByCode(baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	if baseAccount.Balance.Amount < shift.DepositMin {
		log.Fatal("base currency balance is below sideshift minimum")
	}

	if baseAccount.Balance.Amount > shift.DepositMax {
		log.Fatal("base currency balance is above sideshift maximum")
	}

	_, err = cbClient.CreateTransaction(baseAccount.ID, coinbase.TxRequest{
		Type:     "send",
		To:       shift.DepositAddress,
		Amount:   baseAccount.Balance.Amount,
		Currency: baseCurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	shiftResult, err := ssClient.PollShift(shift.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(shiftResult)
}

func createOrder(cbClient *coinbase.Client, currency string) error {
	fiatAccount, err := cbClient.GetAccountByCode("GBP")
	if err != nil {
		return err
	}

	if fiatAccount.Balance.Amount > 0 {
		order := coinbase.CreateAdvancedOrderRequest{
			ClientOrderID: uuid.New().String(),
			ProductID:     fmt.Sprintf("%v-GBP", currency),
			Side:          "BUY",
		}
		order.OrderConfiguration.MarketMarketIOC.QuoteSize = fiatAccount.Balance.Amount

		resp, err := cbClient.CreateAdvancedOrder(order)
		if err != nil {
			return err
		}

		if !resp.Success {
			return errors.New(resp.ErrorResponse.Message)
		}
	}

	return nil
}

func getRefundAddress(cbClient *coinbase.Client, currency string) (string, error) {
	addresses, err := cbClient.GetAddresses(baseCurrency)
	if err != nil {
		return "", err
	}

	var refundAddress string
	if len(*addresses) > 0 {
		refundAddress = (*addresses)[0].Address
	} else {
		address, err := cbClient.CreateAddress(baseCurrency)
		if err != nil {
			return "", err
		}
		refundAddress = address.Address
	}

	return refundAddress, nil
}
