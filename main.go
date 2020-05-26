package main

func main() {
	node := NewNode()
	defer node.cancel()

	crawler := &Crawler{node: node}
	crawler.Start()
}
