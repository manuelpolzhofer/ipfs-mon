package crawler

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"math/rand"
	"sync"
	"time"
)

type Worker struct {
	node        *Node
	maxRoutines int
	mutex       *sync.Mutex
}

func NewWorker(node *Node, maxRoutines int) *Worker {

	return &Worker{node: node, maxRoutines: maxRoutines, mutex: &sync.Mutex{}}
}

func (c *Worker) Start(ctx context.Context, peerCh chan peer.ID, basePeer string, bits int) error {
	fmt.Println("start crawling")
	for i := 0; i < c.maxRoutines; i++ {
		go c.workRoutine(ctx, basePeer, bits, peerCh)
	}

	return nil
}

func (c *Worker) workRoutine(ctx context.Context, basePeer string, bits int, peerCh chan peer.ID) {
	localPeerMap := make(map[string]peer.ID)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p := findPeerWithCommonDHTID(basePeer, bits)
			c.getClosestPeers(ctx, p, basePeer, bits, peerCh, localPeerMap)
		}
	}
}

func (c *Worker) getClosestPeers(ctx context.Context, peerId, basePeer string, bits int, peerCh chan peer.ID, peersMap map[string]peer.ID) {
	// todo figure out optimized tuning for getClosestPeer timeout
	d := rand.Intn(100) + 50
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d)*time.Second)
	defer cancel()

	ch, err := c.node.ipfsNode.DHT.WAN.GetClosestPeers(ctx, peerId)

	if err != nil {
		return
	}

	for {
		select {
		case p := <-ch:
			key := p.String()
			if _, exists := peersMap[key]; !exists {
				peersMap[key] = p
				if kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(string(p))) >= bits {
					peerCh <- p
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
