package cmd

import (
	"fmt"
	"github.com/manuelpolzhofer/ipfs-mon/stats"
	"github.com/spf13/cobra"
)

var (
	fileA string
	fileB string
)

func init() {
	compareCmd.PersistentFlags().StringVar(&fileA, "fileA", "", "path to first crawl result")
	compareCmd.PersistentFlags().StringVar(&fileB, "fileB", "", "path to second crawl result (later in time)")
	rootCmd.AddCommand(compareCmd)
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "compare two result files for overlapping peer IDs",
	Long:  `compare two result files for overlapping peer IDs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ipfs-mon compare...")
		if fileA == "" || fileB == "" {
			fmt.Println("Please set FileA and FileB as param")
		}

		stats.Compare(fileA, fileB)
	},
}
