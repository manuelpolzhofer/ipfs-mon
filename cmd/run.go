package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	basePeerRun  string
	bitsFromBase int
)

func init() {
	runCmd.PersistentFlags().StringVar(&basePeerRun, "basePeer", "", "basePeer for the ipfs node peer id")
	runCmd.PersistentFlags().IntVar(&bitsFromBase, "bitsFromBasePeer", defaultBits, "bits the ipfs peer id should share with the basePeer")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a ipfs node in the same n-bit zone as the given basePeer id",
	Long:  `run a ipfs node in the n-bit zone as the given basePeer id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BasePeer: ", basePeerRun)
		fmt.Println("Bits: ", bitsFromBase)
		fmt.Println("Run IPFS Node in same N-Bit Zone as Base Peer")
		fmt.Println("not implemented yet...")
	},
}
