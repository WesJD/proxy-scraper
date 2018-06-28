package sites

import (
	"github.com/ddliu/go-httpclient"
	"github.com/PuerkitoBio/goquery"
	"time"
	"golang.org/x/net/html"
	"github.com/WesJD/proxy-scraper/utils"
)

type FreeProxyList struct{}

func (s *FreeProxyList) Check(url string) (result map[string]bool, err error) {
	res, err := httpclient.
		Begin().
		Get("https://free-proxy-list.net/anonymous-proxy.html")
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}

	var htmlProxy *html.Node

	result = make(map[string]bool)

	for htmlProxy = doc.Find("#proxylisttable").Get(0).FirstChild.NextSibling.FirstChild;
		htmlProxy.NextSibling != nil;
		htmlProxy = htmlProxy.NextSibling {
		fc := htmlProxy.FirstChild
		address := fc.FirstChild.Data + ":" + fc.NextSibling.FirstChild.Data

		result[address] = utils.CheckProxy(address, url)
	}

	return
}

func (s FreeProxyList) WaitTime() time.Duration {
	return 10 * time.Minute
}
