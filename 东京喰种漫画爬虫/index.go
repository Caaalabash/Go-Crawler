package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"fmt"
	"os"
	"time"
	"strings"
	"io/ioutil"
	"regexp"
	"net/http"
)
//漫画地址
const baseURL  = "https://manhua.fzdm.com/117/"
//是否是图片地址
var isImgURL = regexp.MustCompile(`\.(jpg|png|jpeg)$`)
//并发限制20
var limit = 50

func main(){
	startTime := time.Now()
	downloadLimit := make(chan struct{},limit)

	fmt.Println("开始获取所有章节链接")
	list, e := getChapterLinks()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	fmt.Println(list)
	fmt.Println("开始分析章节链接")

	l := len(list)
	ch := make(chan []string,l)
	for _, href := range list {
		go getPicLinks(href,ch)
	}
	for i:=0;i<l;i++{
		v := <- ch
		fmt.Printf("%s:%d\n",v[len(v)-1],i)
		os.Mkdir(v[len(v)-1],0666)
    for _,imgurl := range v[1:len(v)-1]{
			fmt.Printf("开始下载%s\n",imgurl)
			downloadLimit <- struct{}{}
			go download(imgurl,downloadLimit,v[len(v)-1])
		}
	}
	endTime := time.Now()

	fmt.Println("采集完成,用时",endTime.Sub(startTime))
}

//获取所有章节的链接
func getChapterLinks()([]string,error){
  var list []string
  doc, err := goquery.NewDocument(baseURL)
  if err != nil {
    return nil, fmt.Errorf("fetch %s faild \n: %v",baseURL,err)
  }
  doc.Find(".pure-u-1-2 >a").Each(func(i int,content *goquery.Selection){
    href, e := content.Attr("href")
    if e {
      list = append(list,href)
    }
  })
  return list, nil
}
//获取一个章节下的所有图片链接
func getPicLinks(chap string,ch chan<- []string){
	list := getPages(chap,0,make([]string,1))
	//这里需要返回章节名称 chap格式为 re1123/ 这里需要去除/,然后放在list的最后一位
	chapterTitle := strings.Replace(chap,"/","",-1)
	list = append(list,chapterTitle)
	ch <- list
	return
}

//递归当前章节的页面
func getPages(chap string,index int,list []string) ([]string){
	url := baseURL + chap + "index_" + fmt.Sprint(index) + ".html"
	result := getLinks(url,1)
	//结果为404 说明已经没有下一页了
	if result == "404" {
		return list
	}
	if len(result)>1 {
		list = append(list,result)
	}
	arr := getPages(chap,index+1,list)
	return arr
}

//解析动态js获取图片链接,当请求错误时,有三次重试机会:),在这个阶段就需要尽量减少错误
func getLinks(url string,try int) (result string){
	doc, err := goquery.NewDocument(url)
	//这表明不存在该页 针对该网站的写法
	if err != nil {
		if try<=3 {

			return getLinks(url,try+1)
		}
		return fmt.Sprintf("请求%s错误:尝试次数%d",url,try)
	}
	if doc.Text() == "404"{
		return "404"
	}
	content, err := doc.Find("body > br").Next().Html()
	if err != nil {
		return fmt.Sprintf("寻找动态js内容错误:%v",err)
	}
	//js中的"会变成&#34; 疑惑- - 于是全部替换成单引号
	jsContent := strings.Replace(strings.Split(content,"function")[0],"&#34;","'",-1)
	//图片的地址使用js动态生成,因此捕捉对应的js代码后通过otto库运行js代码获取地址
	vm := otto.New()
	//这里似乎有所变化,原来是101,96.10.31/xx,现在要将其替换为http才可以访问到图片
	vm.Run(`
		function getCookie(val){
			return 'http://p1.xiaoshidi.net'
		}
	`)
	vm.Run(jsContent)
	value, err := vm.Get("mhpicurl")
	if err != nil {
		return fmt.Sprintf("获取%s页面图片请求错误:%v",url,err)
	}
	return fmt.Sprint(value)
}

//下载方法
func download(imgurl string,downloadLimit chan struct{},dirname string){
	//参数检查,imgurl必须是以图片格式结尾
  arr := strings.Split(imgurl,"/")
	name := dirname + arr[len(arr)-1]
	if isImgURL.MatchString(imgurl) == false {
		<- downloadLimit
    return
	}
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("GET", imgurl, nil)
	//增加header选项
	reqest.Header.Add("Referer", "http://manhua.fzdm.com/117/")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Printf("%v\n",err)
		<- downloadLimit
	  return 
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v\n",err)
		<- downloadLimit
		return 
	}
	image, err := os.Create(dirname+"/"+name)
	if err != nil {
		fmt.Printf("%v\n",err)
		<- downloadLimit
		return 
	}
	defer image.Close()
	image.Write(data)
	<- downloadLimit
  return
}