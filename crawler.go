package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
	"math/rand"
	"time"

	kb "github.com/libp2p/go-libp2p-kbucket"
)

type Crawler struct {
	node *Node
}

func (c Crawler) Start() {
	randomPeerId, err := createRandomPeerId()
	if err != nil {
		panic(fmt.Errorf("error random peer id: %s", err))
	}

	fmt.Println("Random Peer Id", randomPeerId)

	ctx, cancel := context.WithCancel(context.Background())

	ch, err := c.node.ipfsNode.DHT.WAN.GetClosestPeers(ctx, randomPeerId)
	if err != nil {
		panic(fmt.Errorf("get closest peers failed: %s", err))
	}

	var peerMap = make(map[string]peer.ID)

	fmt.Println("wait one minute to find more peers")
	time.Sleep(60 * time.Second)
	cancel()

	for peer := range ch {
		if _, exists := peerMap[peer.String()]; !exists {
			peerMap[peer.String()] = peer
		}

	}

	for key, p := range peerMap {
		cp := kb.CommonPrefixLen(kb.ConvertKey(randomPeerId), kb.ConvertKey(string(p)))
		fmt.Println(key, cp)
	}
}

func createRandomPeerId() (string, error) {
	buf := make([]byte, 32)
	rand.Read(buf)

	peerIdBytes, err := mh.Encode(buf, mh.SHA2_256)
	if err != nil {
		return "", err
	}
	return string(peerIdBytes), nil
	//return base58.Encode(peerIdBytes), nil
}
