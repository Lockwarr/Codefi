package scraper

import (
	"log"
	"net/url"

	"golang.org/x/net/html"
)

// CountLinks extracts external & internal links count from a html document
func CountLinks(page *url.URL, document *html.Node) (external, internal uint, err error) {
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" { // only get links from a tags
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "" {
					hrefURL, err := url.Parse(attr.Val)
					if err != nil {
						log.Println("malformed href value: ", attr.Val, err)
						break // nested links are forbidden => assuming there is only one href per node, we can break from the loop
					}
					if (hrefURL.Hostname() == page.Hostname()) ||
						(hrefURL.Hostname() == "" && hrefURL.Path != "") { // if host isn't set but path is set, then the link is most likely internal
						internal++
						continue
					}
					external++
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(document)
	return
}
