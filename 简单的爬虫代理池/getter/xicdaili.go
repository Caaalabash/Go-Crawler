package getter

import (
	"fmt"
	"strings"
	"github.com/parnurzeal/gorequest"
	"github.com/PuerkitoBio/goquery"
)



func Xicdaili() (result []string) {
	pollURL := "http://www.xicidaili.com/nn/1"
	resp, _, errs := gorequest.New().Get(pollURL).
		Set("User-Agent",`"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36"`).
		End()
	if errs != nil {
		fmt.Println(errs)
		return
	}
  doc, e2 := goquery.NewDocumentFromReader(resp.Body)
	if e2 != nil {
		fmt.Printf("error : %v\n",e2)
		return
	}
	doc.Find("#ip_list > tbody > tr").Each(func(i int, s *goquery.Selection) {
		ip := s.Find("td:nth-child(2)").Text()
		port :=s.Find("td:nth-child(3)").Text()
		protocol := strings.ToLower(s.Find("td:nth-child(6)").Text())
    
		data := protocol + "://" + ip + ":" + port
		if i != 0 {
			result = append(result,data)
		}
	})
	fmt.Println("xicdaili done")
	return 
}
