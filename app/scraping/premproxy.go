package scraping

import (
	"time"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"strings"
	"github.com/WesJD/proxy-scraper/app/utils"
	)

type PremProxy struct{}

type HtmlProxy struct {
	IpPort    string
	Anonymity string
}

const (
 	urlFormat = "https://premproxy.com/list/%02d.htm"
 	totalPages = 13 // there are no more than 13 available
)

func (s *PremProxy) Check(url string, trueResponse string) (result map[string]bool, err error) {
	var proxies []string

	result = make(map[string]bool)

	for page := 1; page <= totalPages; page++ {
		proxies, err = getProxies(page)
		if err != nil {
			return
		}
		for _, proxy := range proxies {
			result[proxy] = utils.CheckProxy(url, trueResponse, proxy)
		}
	}

	return
}

func getProxies(pageNumber int) (proxies []string, err error) {
	proxies = make([]string, 0)

	res, err := httpclient.
		Begin().
		Get(fmt.Sprintf(urlFormat, pageNumber))

	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()

	if err != nil {
		return
	}

	for htmlProxyNode := doc.Find("#proxylistt").Nodes[0].FirstChild.NextSibling.NextSibling.NextSibling.FirstChild;
		htmlProxyNode.NextSibling != nil;
		htmlProxyNode = htmlProxyNode.NextSibling {
		proxy := pullProxy(htmlProxyNode)
		if proxy == nil {
			continue
		}
		if proxy.Anonymity != "transparent" {
			proxies = append(proxies, proxy.IpPort)
		}
	}
	return
}

func pullProxy(node *html.Node) *HtmlProxy {
	if node == nil || node.FirstChild == nil {
		return nil
	}
	node = node.FirstChild
	address := strings.TrimSuffix(node.FirstChild.NextSibling.Data, ":")

	if strings.Contains(address, "Select") {
		return nil // non-proxy
	}
	anonymity := node.NextSibling.FirstChild.Data

	return &HtmlProxy{
		IpPort:    address,
		Anonymity: anonymity,
	}
}

func (s *PremProxy) WaitTime() time.Duration {
	return 5 * time.Minute // shrug
}
