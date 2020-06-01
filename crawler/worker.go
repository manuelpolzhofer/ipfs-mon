package crawler

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"sync"
	"time"
)

type Worker struct {
	node        *Node
	maxRoutines int
	mutex       *sync.Mutex
}

const reqSleepTime = 3

func NewWorker(node *Node, maxRoutines int) *Worker {

	return &Worker{node: node, maxRoutines: maxRoutines, mutex: &sync.Mutex{}}
}

func (w *Worker) Start(ctx context.Context, peerCh chan peer.ID, basePeer string, bits int) error {
	fmt.Println("start crawling")
	for i := 0; i < w.maxRoutines; i++ {
		go w.workRoutine(ctx, basePeer, bits, peerCh)
	}

	return nil
}

func (w *Worker) workRoutine(ctx context.Context, basePeer string, bits int, peerCh chan peer.ID) {
	localPeerMap := make(map[string]peer.ID)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p := findPeerWithCommonDHTID(basePeer, bits)
			w.getClosestPeers(ctx, p, basePeer, bits, peerCh, localPeerMap)
		}
	}
}

func (w *Worker) getClosestPeers(ctx context.Context, peerId, basePeer string, bits int, peerCh chan peer.ID, peersMap map[string]peer.ID) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := w.node.ipfsNode.DHT.WAN.GetClosestPeers(ctx, peerId)

	if err != nil {
		panic(fmt.Errorf("get closest peers failed: %s", err))
	}

	time.Sleep(reqSleepTime * time.Second)
	cancel()

	for peer := range ch {
		if _, exists := peersMap[string(peer)]; !exists {
			peersMap[string(peer)] = peer

			if kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(string(peer))) >= bits {
				peerCh <- peer
			}
		}

	}
}
