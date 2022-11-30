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
	CoinbaseToken  string
	SideShiftToken string
	Address        string
	VolumeBase     float64
}

func Convert(opts Opts) {
	ssClient := sideshift.NewClient(opts.SideShiftToken)
	cbClient := coinbase.NewClient("", opts.CoinbaseToken)

	shift, err := ssClient.CreateVariableShift(sideshift.VariableShiftRequest{
		SettleAddress: sideshift.WalletAddress(opts.Address),
		RefundAddress: "",
		DepositCoin:   baseCurrency,
		SettleCoin:    quoteCurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = cbClient.SendTransaction(coinbase.TxRequest{
		Type:     "send",
		To:       shift.DepositAddress,
		Amount:   fmt.Sprintf("%f", opts.VolumeBase),
		Currency: baseCurrency,
		Nonce:    "",
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
