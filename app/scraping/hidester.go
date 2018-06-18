package scraping

import (
	"time"
	"fmt"
	"strconv"
	"github.com/ddliu/go-httpclient"
	"encoding/json"
	"github.com/WesJD/proxy-scraper/app/utils"
)

/*
IP	104.131.22.230
PORT	80
latest_check	1529355004
ping	94
connection_delay	0
country	UNITED STATES
down_speed	0
up_speed	0
proxiescol	null
anonymity	Elite
type	http
google_proxy	0
 */

type Hidester struct {
	Offset int
}

type HidesterProxy struct {
	Ip        string `json:"IP"`
	Port      int    `json:"PORT"`
	Type      string `json:"type"`
	Anonymity string `json:"anonymity"`
}

const checkUrl = "https://hidester.com/proxydata/php/data.php?mykey=data&offset=%s&limit=%s&orderBy=latest_check&sortOrder=DESC&type=http&anonymity=elite&ping=undefined&gproxy=2"
const proxiesPerCheck = 50

func (s *Hidester) Check(url string, trueResponse string) (result map[string]bool, err error) {

	formattedUrl := fmt.Sprintf(checkUrl, strconv.Itoa(s.Offset), strconv.Itoa(proxiesPerCheck))
	result = make(map[string]bool)

	client := httpclient.Begin()

	client.WithHeader("Referer", "https://hidester.com/proxylist/") // bypass their "security"
	res, err := client.Get(formattedUrl)
	if err != nil {
		return
	}

	var response []HidesterProxy

	value, err := res.ToString()
	if err != nil {
		return
	}

	fmt.Println(value)
	if err = json.Unmarshal([]byte(value), &response); err != nil {
		return
	}

	s.Offset++

	result = make(map[string]bool)

	for _, proxy := range response {
		if proxy.Type == "http" && proxy.Anonymity != "Transparent" {
			address := proxy.Ip + ":" + strconv.Itoa(proxy.Port)
			result[address] = utils.CheckProxy(url, trueResponse, address)
		}
	}

	return
}

func (s Hidester) WaitTime() time.Duration {
	return 5 * time.Second
}
