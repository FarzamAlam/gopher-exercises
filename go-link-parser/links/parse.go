package links

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var links []Link
	nodes := linkNodes(doc)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(n *html.Node) Link {
	var link Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}
	link.Text = "to do"
	return link
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func dfs(n *html.Node, padding string) {
	msg := n.Data
	if n.Type == html.ElementNode {
		msg = "<" + msg + ">"
	}
	fmt.Println(padding, msg)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dfs(c, padding+"  ")
	}
}
