package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cedws/fiat2xmr/kraken"
	"github.com/spf13/cobra"
)

var (
	apiKey      string
	privateKey  string
	volumeBase  float32
	volumeQuote float32
)

var rootCmd = &cobra.Command{
	Use: "fiat2xmr",
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := kraken.NewClient(apiKey, privateKey)
		resp, err := client.AddOrder()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", resp)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "a", "", "kraken api key")
	rootCmd.Flags().StringVarP(&privateKey, "private-key", "p", "", "kraken private key")
	rootCmd.Flags().Float32VarP(&volumeBase, "volume-base", "b", 0, "volume to buy using base currency")
	rootCmd.Flags().Float32VarP(&volumeQuote, "volume-quote", "q", 0, "volume to buy using quote currency")

	rootCmd.MarkFlagRequired("api-key")
	rootCmd.MarkFlagRequired("private-key")
	rootCmd.MarkFlagsMutuallyExclusive("volume-base", "volume-quote")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
