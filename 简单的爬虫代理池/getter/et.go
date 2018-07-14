package getter

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type ipInfo struct {
	IP string
	ExpireTime string 
	IpAddress string 
}
type proxyPool struct {
	Data []ipInfo 
}


func EtProxy() (result []string){
	pollURL := "http://47.106.180.108:8081/Index-generate_api_url.html?packid=7&fa=5&qty=10&port=1&format=json&ss=5&css=&ipport=1&pro=&city="
	resp, e := http.Get(pollURL)

	if e != nil {
		fmt.Printf("%v",e)
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Printf("%v",e)
	}
	jsonString := string(body)
  pool := proxyPool{}
	err := json.Unmarshal([]byte(jsonString), &pool)
	if err != nil {
			fmt.Println("error:", err)
	}
	for _,item := range pool.Data{
		result = append(result,"//"+item.IP)
	}
	fmt.Println("et-proxy done")
	return
}