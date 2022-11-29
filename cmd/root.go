package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cedws/fiat2xmr/coinbase"
	"github.com/spf13/cobra"
)

var (
	account    string
	token      string
	address    string
	pair       string
	volumeBase float64
)

var rootCmd = &cobra.Command{
	Use: "fiat2xmr",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := coinbase.NewClient(account, token)
		if err != nil {
			log.Fatal(err)
		}
		_, err = client.SendTransaction(coinbase.TxRequest{
			Type:     "send",
			To:       "",
			Amount:   "0.1",
			Currency: "LTC",
			Nonce:    "123",
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&account, "account", "a", "", "account number")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "account token")
	rootCmd.Flags().StringVarP(&address, "address", "x", "", "monero wallet address")
	rootCmd.Flags().StringVar(&pair, "pair", "LTCXMR", "currency pair to use for conversion")
	rootCmd.Flags().Float64VarP(&volumeBase, "volume-base", "v", 0, "volume to buy using base currency")

	rootCmd.MarkFlagRequired("account")
	rootCmd.MarkFlagRequired("token")
	rootCmd.MarkFlagRequired("address")
	rootCmd.MarkFlagRequired("volume-base")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
