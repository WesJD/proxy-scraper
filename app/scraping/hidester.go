package scraping

import (
	"time"
	"fmt"
	"strconv"
	"github.com/ddliu/go-httpclient"
	"encoding/json"
	"github.com/WesJD/proxy-scraper/app/utils"
	"errors"
)

type Hidester struct {
	Offset int
}

type HidesterProxy struct {
	Ip        string `json:"IP"`
	Port      int    `json:"PORT"`
	Type      string `json:"type"`
	Anonymity string `json:"anonymity"`
}

const (
 	checkUrl = "https://hidester.com/proxydata/php/data.php?mykey=data&offset=%d&limit=%d&orderBy=latest_check&sortOrder=DESC&type=http&anonymity=elite&ping=undefined&gproxy=2"
 	proxiesPerCheck = 50
)

func (s *Hidester) Check(url string, trueResponse string) (result map[string]bool, err error) {
	formattedUrl := fmt.Sprintf(checkUrl, s.Offset, proxiesPerCheck)
	result = make(map[string]bool)

	res, err := httpclient.
		Begin().
		WithHeader("Referer", "https://hidester.com/proxylist/"). // bypass their "security"
		Get(formattedUrl)
	if err != nil {
		return
	}

	var response []HidesterProxy
	value, err := res.ToString()
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(value), &response); err != nil {
		return
	}
	if len(response) == 0 {
		if s.Offset == 0 {
			err = errors.New("no data gotten from offset 0")
			return
		}
		s.Offset = 0
		return s.Check(url, trueResponse)
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
