package crawler

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"

	kb "github.com/libp2p/go-libp2p-kbucket"
)

type Worker struct {
	node        *Node
	maxRoutines int
}

func NewWorker(node *Node, maxRoutines int) *Worker {
	return &Worker{node: node, maxRoutines: maxRoutines}
}

func (c *Worker) Start(ctx context.Context, peerCh chan peer.ID, basePeer string, bits int) error {
	fmt.Println("start crawling")
	for i := 0; i < c.maxRoutines; i++ {
		go c.workRoutine(ctx, basePeer, bits, peerCh)
	}

	return nil
}

func (c *Worker) workRoutine(ctx context.Context, basePeer string, bits int, peerCh chan peer.ID) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop worker")
			return
		default:
			p := findPeerWithCommonDHTID(basePeer, bits)
			c.getClosestPeers(ctx, p, basePeer, bits, peerCh)
		}
	}
}

func (c *Worker) getClosestPeers(ctx context.Context, peerId, basePeer string, bits int, peerCh chan peer.ID) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := c.node.ipfsNode.DHT.WAN.GetClosestPeers(ctx, peerId)
	if err != nil {
		panic(fmt.Errorf("get closest peers failed: %s", err))
	}

	time.Sleep(5 * time.Second)
	cancel()

	for peer := range ch {
		if kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(string(peerId))) >= bits {
			peerCh <- peer
		}
	}
}
