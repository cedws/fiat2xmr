package fiat2xmr

import (
	"fmt"
	"log"
	"time"

	"github.com/cedws/fiat2xmr/coinbase"
	"github.com/cedws/fiat2xmr/sideshift"
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
	DryRun          bool
}

func Convert(opts Opts) {
	cbClient := coinbase.NewClient(opts.CoinbaseKey, opts.CoinbaseSecret)

	account, err := cbClient.GetAccountByCode(baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	refundAddress := getRefundAddress(cbClient, baseCurrency)

	ssClient := sideshift.NewClient(opts.SideShiftSecret)
	shift, err := ssClient.CreateVariableShift(sideshift.VariableShiftRequest{
		SettleAddress: opts.Address,
		RefundAddress: refundAddress,
		DepositCoin:   baseCurrency,
		SettleCoin:    quoteCurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = cbClient.CreateTransaction(account.ID, coinbase.TxRequest{
		Type:     "send",
		To:       shift.DepositAddress,
		Amount:   fmt.Sprintf("%f", opts.VolumeBase),
		Currency: baseCurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	for range time.Tick(10 * time.Second) {
		shift, err := ssClient.GetShift(shift.ID)
		if err != nil {
			log.Fatal(err)
		}

		switch shift.Status {
		case sideshift.StatusSettled:
			break
		case sideshift.StatusWaiting:
		case sideshift.StatusPending:
		case sideshift.StatusProcessing:
		case sideshift.StatusSettling:
			// TODO
		default:
			// TODO
		}
	}
}

func getRefundAddress(cbClient *coinbase.Client, currency string) string {
	addresses, err := cbClient.GetAddresses(baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	var refundAddress string
	if len(*addresses) > 0 {
		refundAddress = (*addresses)[0].Address
	} else {
		address, err := cbClient.CreateAddress(baseCurrency)
		if err != nil {
			log.Fatal(err)
		}
		refundAddress = address.Address
	}

	return refundAddress
}
