package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"

	kb "github.com/libp2p/go-libp2p-kbucket"
)

type Peer struct {
	id           peer.ID
	lastSeen     time.Time
	commonPrefix int
}

type Crawler struct {
	node        *Node
	peersMap    map[string]*Peer
	maxRoutines int
}

func NewCrawler(node *Node) *Crawler {
	m := make(map[string]*Peer)
	return &Crawler{node: node, peersMap: m, maxRoutines: 50}
}

func (c *Crawler) Start(ctx context.Context, peerCh chan peer.ID, basePeer string, bits int) error {
	for i := 0; i < c.maxRoutines; i++ {
		go c.crawlRoutine(ctx, basePeer, bits, peerCh)
	}

	return nil
}

func (c *Crawler) crawlRoutine(ctx context.Context, basePeer string, bits int, peerCh chan peer.ID) {
	fmt.Println("start crawling")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop crawl Routine")
			return
		default:
			p := findPeerWithCommonDHTID(basePeer, bits)
			c.getClosestPeers(ctx, p, basePeer, bits, peerCh)
		}
	}
}

func (c *Crawler) getClosestPeers(ctx context.Context, peerId, basePeer string, bits int, peerCh chan peer.ID) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := c.node.ipfsNode.DHT.WAN.GetClosestPeers(ctx, peerId)
	if err != nil {
		panic(fmt.Errorf("get closest peers failed: %s", err))
	}

	time.Sleep(15 * time.Second)
	cancel()

	for peer := range ch {
		if kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(string(peerId))) >= bits {
			peerCh <- peer
		}
	}
}
