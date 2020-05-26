package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	mh "github.com/multiformats/go-multihash"
	kb "github.com/libp2p/go-libp2p-kbucket"
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

// the Kademlia DHT implementation uses the sha256 of a PID as key
// the method finds a peer ID which has a commonPrefix in the DHT-ID by brute force the sha256-hash function
// estimate: 15 bits ~ 27ms
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
	//return base58.Encode(peerIdBytes), nil
}
