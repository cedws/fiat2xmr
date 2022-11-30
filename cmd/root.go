package cmd

import (
	"fmt"
	"os"

	"github.com/cedws/fiat2xmr/fiat2xmr"
	"github.com/spf13/cobra"
)

var opts fiat2xmr.Opts

var rootCmd = &cobra.Command{
	Use: "fiat2xmr",
	Run: func(cmd *cobra.Command, args []string) {
		fiat2xmr.Convert(opts)
	},
}

func init() {
	rootCmd.Flags().StringVar(&opts.CoinbaseToken, "coinbase-token", "", "coinbase account token")
	rootCmd.Flags().StringVar(&opts.SideShiftToken, "sideshift-token", "", "sideshift account token")
	rootCmd.Flags().StringVarP(&opts.Address, "address", "x", "", "monero wallet address")
	rootCmd.Flags().Float64VarP(&opts.VolumeBase, "volume-base", "v", 0, "volume to buy using base currency")

	rootCmd.MarkFlagRequired("coinbase-token")
	rootCmd.MarkFlagRequired("sideshift-token")
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
