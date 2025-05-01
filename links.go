package main

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	links := []string{}
	r := strings.NewReader(htmlBody)
	nodes, err := html.Parse(r)
	if err != nil {
		return links, err
	}

	for n := range nodes.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link := a.Val
					if string(link[0]) == "/" {
						link = rawBaseURL + link
					}
					links = append(links, link)
					break
				}
			}
		}
	}

	return links, nil
}
