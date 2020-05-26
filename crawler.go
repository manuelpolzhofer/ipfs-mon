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
	return &Crawler{node: node, peersMap: m, maxRoutines: 300}
}

func (c *Crawler) Start() error {
	basePeer, err := createRandomPeerId()
	if err != nil {
		return fmt.Errorf("error random peer id: %s", err)
	}

	basePeerBase58, err := peer.IDFromString(basePeer)
	if err != nil {
		return fmt.Errorf("error convert base peer: %s", err)
	}
	fmt.Println("Base Peer Id", basePeerBase58)

	ctx, cancel := context.WithCancel(context.Background())

	peerCh := make(chan peer.ID, 1000)
	bits := 7
	for i := 0; i < c.maxRoutines; i++ {
		go c.crawlRoutine(ctx, basePeer, bits, peerCh)
	}

	for {
		p := <-peerCh

		if _, exists := c.peersMap[string(p)]; !exists {
			cp := kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(string(p)))
			fmt.Println("New Peer: ", p.String(), "Common Prefix:", cp)
			c.peersMap[string(p)] = &Peer{id: p, lastSeen: time.Now(), commonPrefix: cp}
			fmt.Println("Peers in Zone: ", len(c.peersMap))
		}
		if len(c.peersMap) > 1000 {
			fmt.Println("work done")
			cancel()
			return nil
		}
	}
}

func (c *Crawler) crawlRoutine(ctx context.Context, basePeer string, bits int, peerCh chan peer.ID) {
	for {
		p := findPeerWithCommonDHTID(basePeer, bits)
		c.getClosestPeers(ctx, p, basePeer, bits, peerCh)
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
