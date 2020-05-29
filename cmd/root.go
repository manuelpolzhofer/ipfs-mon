package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ipfs-mon",
	Short: "ipfs-mon is a crawler for the IPFS p2p network",
	Long:  "ipfs-mon is a crawler for the IPFS p2p network.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
