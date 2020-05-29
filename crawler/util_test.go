package crawler

import (
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRandomPeer(t *testing.T) {
	pid, err := createRandomPeerId()
	assert.Nil(t, err)
	p, err := peer.IDFromString(pid)
	assert.Nil(t, err)
	assert.Contains(t, p.String()[:2], "Qm")
}

func TestFindPeerWithCommonDHTID(t *testing.T) {
	for i := 0; i < 100; i++ {
		bits := 15
		basePeer, err := createRandomPeerId()
		assert.Nil(t, err)
		peer := findPeerWithCommonDHTID(basePeer, bits)

		assert.True(t, kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(peer)) >= bits, "common prefix too small")
	}
}
