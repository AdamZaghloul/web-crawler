package main

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error){
	r := strings.NewReader(htmlBody)
	nodes := html.Parse(r)

	for n := range nodes.Descendents(){
		traverseNode(n)
	}
}

func traverseNode(html.Node) (error){
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println(a.Val)
				break
			}
		}
	}
}