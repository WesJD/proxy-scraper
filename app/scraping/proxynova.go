package scraping

import (
	"time"
	"golang.org/x/net/html"
		"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"errors"
	"github.com/WesJD/proxy-scraper/app/utils"
		"fmt"
	"github.com/WesJD/proxy-scraper/app/chrome"
)

type ProxyNova struct{}

type NovaProxy struct {
	Ip        string
	Port      string
	Anonymity string
}

func (s *ProxyNova) Check(url string, trueResponse string) (result map[string]bool, err error) {
	instance, err := chrome.DpInstance("proxynova")
	if err != nil {
		return
	}

	var htmlFull string
	err = instance.Chrome.Run(instance.Context, chromedp.Tasks{
		chromedp.Navigate("https://www.proxynova.com/proxy-server-list/"),
		chromedp.WaitVisible("#tbl_proxy_list"),
		chromedp.OuterHTML("html", &htmlFull),
	})
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlFull))
	if err != nil {
		return
	}

	var htmlProxy *html.Node
	proxyList := doc.Find("#tbl_proxy_list")

	if len(proxyList.Nodes) < 1 {
		err = errors.New("unable to find table")
		return
	}

	result = make(map[string]bool)

	for htmlProxy = proxyList.Nodes[0].FirstChild.NextSibling.NextSibling.NextSibling.FirstChild;
		htmlProxy != nil;
		htmlProxy = htmlProxy.NextSibling {
		if htmlProxy.Data != "tr" {
			continue // some blank space shit
		}
		proxy := parseProxy(htmlProxy)
		if proxy == nil || proxy.Anonymity == "Transparent" {
			continue
		}
		address := proxy.Ip + ":" + proxy.Port
		result[address] = utils.CheckProxy(url, trueResponse, address)
	}

	fmt.Println(":(")

	return
}

func parseProxy(proxy *html.Node) *NovaProxy {
	parse := proxy.FirstChild.NextSibling
	if parse == nil {
		return nil
	}

	address := strings.TrimSuffix(parse.FirstChild.NextSibling.FirstChild.NextSibling.Data, " ")
	parse = parse.NextSibling.NextSibling

	port := clean(parse.FirstChild.Data)
	if len(port) < 1 {
		port = clean(parse.FirstChild.NextSibling.FirstChild.Data) // the port is a link
	}

	parse = parse.NextSibling.NextSibling.
		NextSibling.NextSibling.
		NextSibling.NextSibling.
		NextSibling.NextSibling.
		NextSibling.NextSibling

	anonymity := parse.FirstChild.NextSibling.FirstChild.Data

	return &NovaProxy{
		Ip:        address,
		Port:      port,
		Anonymity: anonymity,
	}
}

func clean(input string) string {
	var replace = []string{
		" ",
		"\t",
		"\n",
	}
	for _, replacing := range replace {
		input = strings.Replace(input, replacing, "", -1)
	}
	return input
}

func (s *ProxyNova) WaitTime() time.Duration {
	return 5 * time.Minute
}
