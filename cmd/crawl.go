package cmd

import (
	"context"
	"fmt"
	"github.com/manuelpolzhofer/ipfs-mon/crawler"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var (
	bits     int
	workers  int
	nodes    int
	maxPeers int
)

const (
	defaultBits     = 6
	defaultWorker   = 100
	defaultNodes    = 1
	defaultMaxPeers = 10000
)

func init() {
	crawlCmd.PersistentFlags().IntVar(&bits, "bits", defaultBits, "define the n-bit zone for a crawl. (0-256)")
	crawlCmd.PersistentFlags().IntVar(&workers, "workers", defaultWorker, "define the amount of workers per node for the crawl")
	crawlCmd.PersistentFlags().IntVar(&nodes, "nodes", defaultNodes, "define the amount of ipfs nodes for the crawl")
	crawlCmd.PersistentFlags().IntVar(&maxPeers, "maxPeers", defaultMaxPeers, "stop crawling after maxPeers are found in zone")
	rootCmd.AddCommand(crawlCmd)
}

func shutdownHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Starting Shutdown")
		cancel()
	}()
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "crawl for other peers which share a common n-bit Kademlia ID",
	Long:  `crawl for other peers which share a common n-bit Kademlia ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing ipfs-mon...")
		fmt.Println(fmt.Sprintf("N-Bit-Zone: %d Bits", bits))
		fmt.Println(fmt.Sprintf("Nodes: %d ", nodes))
		fmt.Println(fmt.Sprintf("Workers per Node: %d ", workers))
		ctx, cancel := context.WithCancel(context.Background())
		shutdownHandler(cancel)
		c := crawler.NewCluster(defaultNodes, workers, bits, maxPeers)
		c.Start(ctx)
	},
}
