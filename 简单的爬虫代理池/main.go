package main

import (
	"fmt"
	"time"
	"sync"
	"./getter"
	"github.com/parnurzeal/gorequest"
)

func main (){
	start:= time.Now()
	ipChan := make(chan string, 200)
	ipList := make([]string,200)
	proxyList := make(chan string,20)
  
	go Run(ipChan)

	for ip := range ipChan{
		ipList = append(ipList,ip)
	}
  fmt.Printf("共采集 %d 条ip\n",len(ipList))

	go RunCheckIp(ipList,proxyList)

  for v := range proxyList{
		fmt.Println(v,"可用")
	}
	end := time.Now()
  fmt.Println("程序运行结束",end.Sub(start))
}

//爬取getter中的代理网站,获得ip
func Run(ipChan chan<-string){
	var wg sync.WaitGroup
  funs := []func() []string{
		getter.Data5u,
		getter.EtProxy,
		getter.IP66,
		getter.Xicdaili,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() []string) {
			temp := f()
			for _, v := range temp {
				ipChan <- v
			}
			defer wg.Done()
		}(f)
	}
	wg.Wait()
	fmt.Println("执行run完成")
	close(ipChan)
}

func RunCheckIp(list []string,proxyList chan string) {
	var wg sync.WaitGroup
	for _,ip := range list {
		wg.Add(1)
		go func(ip string){
			flag := checkIp(ip)
			if flag==true{
				proxyList <- ip
			}
			defer wg.Done()
		}(ip)
	}
	wg.Wait()
	fmt.Println("执行check完成")
	close(proxyList)
}


//检测代理是否可用
func checkIp(ip string) bool {
	if len(ip) <2 {
		return false
	}
	pollURL := "https://blog.calabash.top"
	resp, _, errs := gorequest.New().Timeout(10*time.Second).Proxy(ip).Get(pollURL).
		Set("User-Agent",`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36"`).
		End()
	if errs != nil {
		return false
	}
	if resp.StatusCode == 200 {
	  return true
	}
	return false
}