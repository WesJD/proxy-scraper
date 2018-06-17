package scraping

import (
	"github.com/ddliu/go-httpclient"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/PuerkitoBio/goquery"
)

func Show() {
	res, err := httpclient.Begin().Get("url")
	utils.CheckError(err)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckError(err)
	//use doc
}
