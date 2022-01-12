package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func ScanNode(node *html.Node) *html.Node {
	if node.Data == "blockroot" {
		return node
	}
	for _node := node.FirstChild; _node != nil; _node = _node.NextSibling {
		if _n := ScanNode(_node); _n != nil {
			return _n
		}
	}
	return nil
}

func main() {
	doc, _ := html.Parse(strings.NewReader(htm))
	node := ScanNode(doc)
	fmt.Println(node.Data, node.Attr)
}

const htm = `
<asdf />
<a>
<blockroot name="root" kind="basic" pivot="0,0,0,0" size="0,0,0,0" pos="0,0,0,0" rotation="0,0,0,0">
	<block name="body">
	</block>
</blockroot>
</a>
`
