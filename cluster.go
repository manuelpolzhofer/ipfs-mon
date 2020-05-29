package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"os"
	"time"
)

type Cluster struct {
	peersMap map[string]*Peer
	peerCh   chan peer.ID
	bits     int
	maxPeers int
	numNodes int
	basePeer string
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewCluster() *Cluster {
	m := make(map[string]*Peer)
	return &Cluster{peersMap: m, bits: 6, numNodes: 1, maxPeers: 10000}
}

func (c *Cluster) Start(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	c.peerCh = make(chan peer.ID, 1000)

	node := NewNode(ctx, "")
	defer node.cancel()
	crawler := NewCrawler(node)

	// use peerID of first node as basePeer
	// the first n bits of the basePeer are used for finding other peers in the zone
	c.basePeer = string(node.ipfsNode.Identity)

	crawler.Start(ctx, c.peerCh, c.basePeer, c.bits)

	for i := 0; i < c.numNodes-1; i++ {
		go func(ctx context.Context, peerCh chan peer.ID, basePeer string, bits int) {

			node := NewNode(ctx, basePeer)
			defer node.cancel()
			crawler := NewCrawler(node)

			crawler.Start(ctx, peerCh, basePeer, bits)
			<-ctx.Done()
		}(c.ctx, c.peerCh, c.basePeer, c.bits)
	}

	c.listenForPeers()

	return nil
}

// listen for peerID from the crawlers
// each peerID is in the n-bit zone
func (c *Cluster) listenForPeers() {
	fmt.Println("Start listening for peers in zone")
	for {
		select {
		case p := <-c.peerCh:
			c.handleNewPeer(p)

		case <-c.ctx.Done():
			c.handleShutdown()
			return
		}
	}
}

func (c *Cluster) handleNewPeer(p peer.ID) {
	key := p.String()
	if _, exists := c.peersMap[key]; !exists {
		cp := kb.CommonPrefixLen(kb.ConvertKey(c.basePeer), kb.ConvertKey(string(p)))
		c.peersMap[key] = NewPeer(p, time.Now(), cp)

		fmt.Println("New Peer: ", p.String(), "Common Prefix:", cp)
		fmt.Println("Peers in Zone: ", len(c.peersMap))

	}
	if len(c.peersMap) > c.maxPeers {
		fmt.Println("Finish. Reached Max Peers in Zone: ", len(c.peersMap))
		c.cancel()
		return
	}
}

func (c *Cluster) handleShutdown() {
	fmt.Println("Create JSON file")
	b, err := json.Marshal(c)
	if err != nil {
		panic(fmt.Errorf("failed to create json: %s", err))
	}

	writeResultJSON(b)
	os.Exit(0)
}

func (c *Cluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		InZone    int              `json:"total_peers_in_zone"`
		Bits      int              `json:"bits"`
		BasePeer  string           `json:"base_peer"`
		BaseKADId string           `json:"base_kad_id"`
		PeersMap  map[string]*Peer `json:"peers"`
	}{
		InZone:    len(c.peersMap),
		Bits:      c.bits,
		BasePeer:  peerIDtoBase58(c.basePeer),
		BaseKADId: "0x" + hex.EncodeToString(kb.ConvertKey(c.basePeer)),
		PeersMap:  c.peersMap,
	})
}
