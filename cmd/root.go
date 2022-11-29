package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiKey      string
	privateKey  string
	pair        string
	volumeBase  float64
	volumeQuote float64
)

var rootCmd = &cobra.Command{
	Use: "fiat2xmr",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "a", "", "kraken api key")
	rootCmd.Flags().StringVarP(&privateKey, "private-key", "p", "", "kraken private key")
	rootCmd.Flags().StringVar(&pair, "pair", "LTCXMR", "currency pair to use for conversion")
	rootCmd.Flags().Float64VarP(&volumeBase, "volume-base", "b", 0, "volume to buy using base currency")
	rootCmd.Flags().Float64VarP(&volumeQuote, "volume-quote", "q", 0, "volume to buy using quote currency")

	rootCmd.MarkFlagRequired("api-key")
	rootCmd.MarkFlagRequired("private-key")
	rootCmd.MarkFlagRequired("volume-base")
	rootCmd.MarkFlagsMutuallyExclusive("volume-base", "volume-quote")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
