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
	rootCmd.Flags().StringVar(&opts.CoinbaseKey, "coinbase-key", "", "coinbase account key")
	rootCmd.Flags().StringVar(&opts.CoinbaseSecret, "coinbase-secret", "", "coinbase account secret")
	rootCmd.Flags().StringVar(&opts.SideShiftSecret, "sideshift-secret", "", "sideshift account secret")
	rootCmd.Flags().StringVarP(&opts.Address, "address", "x", "", "monero wallet address")
	rootCmd.Flags().Float64VarP(&opts.VolumeBase, "volume-base", "v", 0, "volume to buy using base currency")

	rootCmd.MarkFlagRequired("coinbase-key")
	rootCmd.MarkFlagRequired("coinbase-secret")
	rootCmd.MarkFlagRequired("sideshift-secret")
	rootCmd.MarkFlagRequired("address")
	rootCmd.MarkFlagRequired("volume-base")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
