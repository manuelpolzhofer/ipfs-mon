# ipfs-mon
A crawler for IPFS nodes to estimate the number of nodes and daily churn in IPFS


### Crawling n-bit zones in the Kademlia DHT 
The crawler is only crawling one specific predefined n-bit zone (sub tree) of a the Kademlia ID space.

The amount of total peers found in a specific n-bit zone can be used to calculate a probabilistic estimate of the
total nodes in the p2p network. 

#### Paper 
The idea for the crawler is inspired by the following paper

 [Measuring Large-Scale Distributed Systems: Case of BitTorrent Mainline DHT](https://www.cs.helsinki.fi/u/lxwang/publications/P2P2013_13.pdf)


## Build
```
go build
```

## Usage
### Help
```bash
./ipfs-mon --help
```

### Crawl Command
Starts a crawl of the IPFS p2p network

#### Example
```bash
./ipfs-mon crawl --bits=5 --workers=300 --basePeer=QmbNCbBpuRCCPBEqvpKNT77ngbEqnGfDbPR4YD7HvroU9C
```

For a `crawl` the following parameters can be provided

| Param | Desc | 
| -------- | -------- |
| --basePeer  | base58 peerID. the first n-bits of the basePeer are used for the zone. (default: ipfs node peerID)     | |
| --bits    | minimum amount of bits a found peer needs to have in common with the basePeer|
| --nodes    | total amount of nodes a crawler should spin up (currently only support for one node) |
| --workers    | total amount of workers per node|
| --peersFile    | file with additional peers for bootstrapping |
|--maxPeers|  if the maxPeers amount is reached the crawler stops. (default: 100000) |

### Compare Command
Compare different results JSON files.

#### Example
```bash
./ipfs-mon compare --fileA=ipfs-mon-crawl-1590869793.json --fileB=ipfs-mon-crawl-1590957682.json
```

For a `compare` the following parameters can be provided

| Param | Desc | 
| -------- | -------- |
| --basePeer  | the ipfs node uses the first n-bits of the basePeer and runs in the same zone| |
| --bits    | minimum amount of bits the ipfs nodes needs to have in common with the base peer|


