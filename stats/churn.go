package stats

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"github.com/manuelpolzhofer/ipfs-mon/zone"
	"io/ioutil"
	"time"
)

type Result struct {
	InZone    int                   `json:"total_peers_in_zone"`
	Bits      int                   `json:"bits"`
	BasePeer  string                `json:"base_peer"`
	BaseKADId string                `json:"base_kad_id"`
	PeersMap  map[string]*zone.Peer `json:"peers"`
	StartTime time.Time             `json:"start_time"`
	EndTime   time.Time             `json:"end_time"`
}

func Compare(fileA, fileB string) {
	resultA := parseResult(fileA)
	resultB := parseResult(fileB)

	fmt.Println("\nResult Summary File A:")
	statsOverview(resultA)

	fmt.Println("\nResult Summary File B:")
	statsOverview(resultB)

	basePeerA, _ := peer.Decode(resultA.BasePeer)
	basePeerB, _ := peer.Decode(resultB.BasePeer)

	cp := kb.CommonPrefixLen(kb.ConvertKey(string(basePeerA)), kb.ConvertKey(string(basePeerB)))

	fmt.Println("\nCompare both Files")
	fmt.Println(">>> Both BasePeer Common Prefix:", cp)
	interSec := peerIntersection(resultA.PeersMap, resultB.PeersMap)
	lenA := len(resultA.PeersMap)
	lenB := len(resultB.PeersMap)
	fmt.Println(">>> Total Peers File A: ", lenA)
	fmt.Println(">>> Total Peers File B: ", lenB)
	fmt.Println(">>> Intersection of Peers: ", interSec)
	fmt.Println(fmt.Sprintf(">>> Peers from file A not in fileB:  %d", lenA-interSec))
	fmt.Println(fmt.Sprintf(">>> Peers from file B not in fileA:  %d", lenB-interSec))
	fmt.Println(">>> Time difference between two crawls: ", resultB.StartTime.Sub(resultA.EndTime))

}

func statsOverview(result *Result) {
	fmt.Println(">>> BasePeer:", result.BasePeer)
	fmt.Println(">>> Peers Total: ", len(result.PeersMap))
	fmt.Println(">>> Bit Zone: ", result.Bits)
	fmt.Println(">>> Start Time:", result.StartTime)
	fmt.Println(">>> End Time:", result.EndTime)
	fmt.Println(">>> Total Time:", result.EndTime.Sub(result.StartTime))
}

func peerIntersection(peersMapA, peersMapB map[string]*zone.Peer) int {
	interSec := 0
	for peerID, _ := range peersMapA {
		if _, exists := peersMapB[peerID]; exists {
			interSec++
		}
	}
	return interSec
}

func parseResult(file string) *Result {
	jsonFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %s", err))
	}

	crawlA := Result{}

	err = json.Unmarshal([]byte(jsonFile), &crawlA)
	if err != nil {
		panic(fmt.Errorf("failed to parse json file: %s", err))
	}

	return &crawlA

}
