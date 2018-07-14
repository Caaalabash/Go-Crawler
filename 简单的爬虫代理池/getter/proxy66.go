package getter

import (
	"fmt"
	"regexp"
	"github.com/parnurzeal/gorequest"
)


// IP66 get ip from 66ip.cn
func IP66() (result []string) {
	reg := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\:\d{1,5}`
	re := regexp.MustCompile(reg)
	pollURL := "http://www.66ip.cn/mo.php?tqsl=100"
	_, body, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		fmt.Println(errs)
		return
	}
	list := re.FindAllString(string(body),-1)
	for _,item := range list{
		result = append(result,"//"+item)
	}
	fmt.Println("66proxy done")
	return 
}
