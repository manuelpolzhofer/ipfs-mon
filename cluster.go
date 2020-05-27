package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"time"
)

type Cluster struct {
	peersMap map[string]*Peer
	peerCh   chan peer.ID
	bits     int
	n        int
	numNodes int
	basePeer string
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewCluster() *Cluster {
	m := make(map[string]*Peer)
	return &Cluster{peersMap: m, bits: 7, numNodes: 1, n: 200}
}

func (c *Cluster) Start() error {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	var err error
	c.basePeer, err = createRandomPeerId()
	if err != nil {
		return fmt.Errorf("error random peer id: %s", err)
	}
	fmt.Println("Base Peer Id", peerIDtoBase58(c.basePeer))
	fmt.Println("Common Bits for Zone: ", c.bits)

	c.peerCh = make(chan peer.ID, 1000)

	for i := 0; i < c.numNodes; i++ {

		go startNode(c.ctx, c.peerCh, c.basePeer, c.bits)
	}

	c.listenForPeers()

	time.Sleep(10 * time.Second)
	return nil
}

func startNode(ctx context.Context, peerCh chan peer.ID, basePeer string, bits int) {
	node := NewNode(ctx)
	defer node.cancel()

	crawler := NewCrawler(node)
	crawler.Start(ctx, peerCh, basePeer, bits)
	<-ctx.Done()
}

// listen for peerID from the crawlers
// each peerID is in the n-bit zone
func (c *Cluster) listenForPeers() {
	fmt.Println("Start listening for peers in zone")
	for {
		p := <-c.peerCh

		if _, exists := c.peersMap[string(p)]; !exists {
			cp := kb.CommonPrefixLen(kb.ConvertKey(c.basePeer), kb.ConvertKey(string(p)))
			c.peersMap[string(p)] = &Peer{id: p, lastSeen: time.Now(), commonPrefix: cp}
			fmt.Println("New Peer: ", p.String(), "Common Prefix:", cp)
			fmt.Println("Peers in Zone: ", len(c.peersMap))

		}
		if len(c.peersMap) > c.n {
			fmt.Println("Finish. Total peers in zone: ", len(c.peersMap))
			c.cancel()
			return
		}
	}
}
