package zone

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

type Peer struct {
	id           peer.ID
	lastSeen     time.Time
	commonPrefix int
}

func NewPeer(id peer.ID, lastSeen time.Time, commonPrefix int) *Peer {
	return &Peer{id: id, lastSeen: lastSeen, commonPrefix: commonPrefix}
}

func (p *Peer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id           string    `json:"id"`
		LastSeen     time.Time `json:"last_seen"`
		CommonPrefix int       `json:"common_prefix"`
	}{
		Id:           p.id.String(),
		LastSeen:     p.lastSeen,
		CommonPrefix: p.commonPrefix,
	})
}
