package crawler

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	mh "github.com/multiformats/go-multihash"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func ReadFile(path string) ([]string, error) {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not read peers file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read peers file: %s", err)
	}
	return lines, nil
}

func writeResultJSON(b []byte) {
	fileName := "ipfs-mon-crawl-" + strconv.Itoa(int(time.Now().Unix())) + ".json"

	fmt.Println("Generated file with results: ", fileName)
	err := ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		panic(fmt.Errorf("error convert base peer: %s", err))
	}

}

// the Kademlia DHT implementation uses the sha256 of a PID as key
// the method finds a peer ID which shares a common prefix in the DHT-ID by brute force the sha256-hash function
func findPeerWithCommonDHTID(basePeer string, bits int) string {
	p, err := createRandomPeerId()
	if err != nil {
		panic(fmt.Errorf("failed to create peer id with common prefix: %s", err))
	}
	for kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(p)) < bits {
		p, err = createRandomPeerId()
		if err != nil {
			panic(fmt.Errorf("failed to create peer id with common prefix: %s", err))
		}
	}
	return p
}

func createRandomPeerId() (string, error) {
	buf := make([]byte, 32)
	rand.Read(buf)

	peerIdBytes, err := mh.Encode(buf, mh.SHA2_256)
	if err != nil {
		return "", err
	}
	return string(peerIdBytes), nil
}

func toBase58(peerID string) string {
	id, err := peer.IDFromString(peerID)
	if err != nil {
		panic(fmt.Errorf("error convert base peer: %s", err))
	}
	return id.Pretty()
}

func fromBase58(peerID string) string {
	id, err := peer.Decode(peerID)
	if err != nil {
		panic(fmt.Errorf("basePeer Id incorrect: %s", err))
	}
	return string(id)
}
