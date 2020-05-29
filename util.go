package main

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	kb "github.com/libp2p/go-libp2p-kbucket"
	mh "github.com/multiformats/go-multihash"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func ReadFile(path string) []string {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func writeResultJSON(b []byte) {
	fileName := "ipfs-mon-crawl-" + strconv.Itoa(int(time.Now().Unix())) + ".json"

	fmt.Println("generated file with results: ", fileName)
	err := ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		panic(fmt.Errorf("error convert base peer: %s", err))
	}

}

// the Kademlia DHT implementation uses the sha256 of a PID as key
// the method finds a peer ID which shares a common prefix in the DHT-ID by brute force the sha256-hash function
func findPeerWithCommonDHTID(basePeer string, bits int) string {
	var p string
	for kb.CommonPrefixLen(kb.ConvertKey(basePeer), kb.ConvertKey(p)) < bits {
		p, _ = createRandomPeerId()

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

func peerIDtoBase58(peerID string) string {
	peerBase58, err := peer.IDFromString(peerID)
	if err != nil {
		panic(fmt.Errorf("error convert base peer: %s", err))
	}
	return peerBase58.String()
}
