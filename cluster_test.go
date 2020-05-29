package main

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestJSONMarshal(t *testing.T) {
	m := make(map[string]*Peer)

	p, err := createRandomPeerId()
	assert.Nil(t, err)
	peer, err := peer.IDFromString(p)
	assert.Nil(t, err)

	m[peer.String()] = &Peer{id: peer, lastSeen: time.Now(), commonPrefix: 3}
	c := &Cluster{peersMap: m, bits: 3, numNodes: 1, maxPeers: 10000, basePeer: p}

	s, err := json.Marshal(c)
	assert.Nil(t, err)

	assert.Contains(t, string(s), peer.String())
}
