package cmd

import (
	"log"

	"github.com/cedws/fiat2xmr/coinbase"
	"github.com/cedws/fiat2xmr/fiat2xmr"
	"github.com/cedws/fiat2xmr/sideshift"
	"github.com/spf13/cobra"
)

var opts fiat2xmr.Opts

var rootCmd = &cobra.Command{
	Use: "fiat2xmr",
	Run: func(cmd *cobra.Command, args []string) {
		ssClient := sideshift.NewClient(opts.SideShiftSecret)
		if canShift, err := ssClient.CanShift(); !canShift || err != nil {
			if !canShift {
				log.Fatal("sideshift account is unable to create shifts")
			}
			log.Fatal(err)
		}
		cbClient := coinbase.NewClient(opts.CoinbaseKey, opts.CoinbaseSecret)

		cnv := fiat2xmr.NewConverter(ssClient, cbClient)
		cnv.Convert(opts)
	},
}

func init() {
	rootCmd.Flags().StringVar(&opts.CoinbaseKey, "coinbase-key", "", "coinbase account key")
	rootCmd.Flags().StringVar(&opts.CoinbaseSecret, "coinbase-secret", "", "coinbase account secret")
	rootCmd.Flags().StringVar(&opts.SideShiftSecret, "sideshift-secret", "", "sideshift account secret")
	rootCmd.Flags().StringVarP(&opts.Address, "address", "x", "", "monero wallet address")

	rootCmd.MarkFlagRequired("coinbase-key")
	rootCmd.MarkFlagRequired("coinbase-secret")
	rootCmd.MarkFlagRequired("sideshift-secret")
	rootCmd.MarkFlagRequired("address")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
