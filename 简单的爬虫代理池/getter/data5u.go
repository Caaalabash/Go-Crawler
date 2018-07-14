package getter

import (
	"fmt"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

func Data5u () (result []string) {
	pollURL := "http://www.data5u.com/free/gngn/index.shtml"
	resp, _, e1 := gorequest.New().Get(pollURL).
		Set("User-Agent",`"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36"`).
		End()
  if e1 != nil {
		fmt.Printf("error : %v\n",e1)
		return
	}
	doc, e2 := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	if e2 != nil {
		fmt.Printf("error : %v\n",e2)
		return
	}
	doc.Find("div.wlist > ul > li:nth-child(2) > ul").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		ip := s.Find("ul:nth-child(" + node + ") > span:nth-child(1) > li").Text()
		port := s.Find("ul:nth-child(" + node + ") > span:nth-child(2) > li").Text()
		protocol := s.Find("ul:nth-child(" + node + ") > span:nth-child(4) > li").Text()

		data := protocol + "://" + ip + ":" + port
		if i != 0{
			result = append(result,data)
		}
	})
	fmt.Println("data5u done")
	return
}