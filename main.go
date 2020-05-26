package main

func main() {
	node := NewNode()
	defer node.cancel()

	crawler := NewCrawler(node)
	crawler.Start()
}
