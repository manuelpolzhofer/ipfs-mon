package crawler

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
	peersMap       map[string]*Peer
	peerCh         chan peer.ID
	bits           int
	maxPeers       int
	numNodes       int
	basePeer       string
	workers        int
	peersFile      string
	ctx            context.Context
	cancel         context.CancelFunc
	firstPeerFound time.Time
	startTime      time.Time
	endTime        time.Time
	lastNewPeer    time.Time
}

func NewCluster(numNodes, workers, bits, maxPeers int, basePeer, peersFile string) *Cluster {
	m := make(map[string]*Peer)

	if basePeer != "" {
		peerID, err := peer.Decode(basePeer)
		if err != nil {
			panic(fmt.Errorf("basePeer Id incorrect: %s", err))
		}
		basePeer = string(peerID)
	}

	return &Cluster{peersMap: m, numNodes: numNodes, bits: bits, workers: workers, maxPeers: maxPeers, basePeer: basePeer, peersFile: peersFile}
}

func (c *Cluster) Start(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	c.peerCh = make(chan peer.ID, 1000)

	c.startTime = time.Now()

	node := NewNode(ctx, c.basePeer, c.peersFile)

	// use peerID of first node as basePeer if not set
	// the first n bits of the basePeer are used for finding other peers in the zone
	if c.basePeer == "" {
		c.basePeer = string(node.ipfsNode.Identity)
	}

	fmt.Println("BasePeer:", peerIDtoBase58(c.basePeer))
	fmt.Println("IPFS Node Peer ID:", peer.Encode(node.ipfsNode.Identity))

	defer node.cancel()
	worker := NewWorker(node, c.workers)

	worker.Start(ctx, c.peerCh, c.basePeer, c.bits)

	for i := 0; i < c.numNodes-1; i++ {
		go func(ctx context.Context, peerCh chan peer.ID, basePeer string, bits, workers int) {

			node := NewNode(ctx, basePeer, c.peersFile)
			defer node.cancel()
			worker := NewWorker(node, workers)

			worker.Start(ctx, peerCh, basePeer, bits)
			<-ctx.Done()
		}(c.ctx, c.peerCh, c.basePeer, c.bits, c.workers)
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

		if len(c.peersMap) == 1 {
			c.firstPeerFound = time.Now()
		}

		fmt.Println("New Peer: ", p.String(), "Common Prefix:", cp, "since:", time.Since(c.lastNewPeer))
		fmt.Println("Peers in Zone: ", len(c.peersMap))
		c.lastNewPeer = time.Now()

	}
	if len(c.peersMap) > c.maxPeers {
		fmt.Println("Finish. Reached Max Peers in Zone: ", len(c.peersMap))
		c.cancel()
		c.handleShutdown()
		return
	}
}

func (c *Cluster) handleShutdown() {
	c.endTime = time.Now()
	totalTime := time.Since(c.startTime)
	fmt.Println("Time past since First Peer Found: ", time.Since(c.firstPeerFound))
	fmt.Println("Time past since last New Peer: ", time.Since(c.lastNewPeer))
	fmt.Println("Total Crawl Time:", totalTime)
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
