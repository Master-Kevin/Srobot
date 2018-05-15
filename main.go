package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// regHref       = `((ht|f)tps?)://[w]{0,3}.baidu.com/link\?[a-zA-z=0-9-\s]*`
	regTitle      = `<title[\sa-zA-z="-]*>([^x00-xff]|[\sa-zA-Z=-：|，？"])*</title>`
	regCheckTitle = `(为什么|怎么)*.*([G|g][O|o][L|l][A|a][N|n][G|g]).*(怎么|如何|为什么).*`
)

func main() {
	if checkFile("./data/", "url.txt").Size() == 0 {
		fistStart()
		main()
	} else {
		Timer()
	}
}

func Timer() {
	t := time.NewTimer(time.Second * 1)
	<-t.C
	fmt.Print("\n\n\n执行爬抓\n\n")
	f, _ := os.OpenFile("./data/url.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	file, _ := ioutil.ReadAll(f)
	pageCont, _ := pageVisit(strings.Split(string(file), "\n")[0])
	if checkRegexp(checkRegexp(pageCont, regTitle, 0).(string), regCheckTitle, 0).(string) != "" {
		fmt.Print(checkRegexp(checkRegexp(pageCont, regTitle, 0).(string), regCheckTitle, 0).(string))
		fmt.Print("\n有效内容 => " + checkRegexp(pageCont, regTitle, 0).(string))
	}
	fmt.Print("\n\n待爬抓网址共" + strconv.Itoa(len(strings.Split(string(file), "\n"))-1) + "个 => " + strings.Split(string(file), "\n")[0] + "\n")
	// DelFirstText("./data/url.txt")
	Timer()
}

func fistStart() {
	var num int
	url := "http://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&tn=39042058_20_oem_dg&wd=golang%E5%AE%9E%E7%8E%B0&oq=golang%2520%25E5%2588%25A0%25E9%2599%25A4%25E6%2595%25B0%25E7%25BB%2584&rsv_pq=d9be28ec0002df1b&rsv_t=8017GWpSLPhDmKilZQ1StC04EVpUAeLEP90NIm%2Bk5pRh5R9o57NHMO8Gaxm1TtSOo%2FvtJj%2B98%2Fsc&rqlang=cn&rsv_enter=1&inputT=3474&rsv_sug3=16&rsv_sug1=11&rsv_sug7=100&rsv_sug2=0&rsv_sug4=4230"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`((ht|f)tps?)://[w]{0,3}.baidu.com/link\?[a-zA-z=0-9-\s]*`)
	f, _ := os.OpenFile("./data/url.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	for _, d := range reg.FindAllString(string(body), -1) {
		ff, _ := os.OpenFile("./data/url.txt", os.O_RDWR, 0666)
		file, _ := ioutil.ReadAll(ff)
		dd := strings.Split(d, "")
		dddd := ""
		for _, ddd := range dd {
			if ddd == "?" {
				ddd = `\?`
			}
			dddd += ddd
		}
		if checkRegexp(string(file), dddd, 0).(string) == "" {
			io.WriteString(f, d+"\n")
			fmt.Print("\n收集地址：" + d + "\n")
			num++
		}
		// fmt.Print(string(file))
		ff.Close()
	}
	fmt.Print("\n首次收集网络地址：" + strconv.Itoa(len(reg.FindAllString(string(body), -1))) + "\n")
	fmt.Print("\n去重后网络地址数：" + strconv.Itoa(num))
	fmt.Print("\n\n首次储存成功！\n")
}

func pageVisit(url string) (page string, body []byte) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	page = string(body)
	return
}

func checkFile(dir string, file string) os.FileInfo {
	list, _ := ioutil.ReadDir(dir)
	for _, info := range list {
		if info.Name() == file {
			return info
		}
	}
	return list[0]
}

func saveFile(file string, cont string) {
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	io.WriteString(f, cont)
}

func checkRegexp(cont string, reg string, style int) (result interface{}) {
	check := regexp.MustCompile(reg)
	switch style {
	case 0:
		result = check.FindString(cont)
	case 1:
		result = check.FindAllString(cont, -1)
	default:
		result = check.FindAll([]byte(cont), -1)
	}
	return
}

func DelFirstText(file string) {
	var text = ""
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	files, _ := ioutil.ReadAll(f)
	var ss = strings.Split(string(files), "\n")
	for i := 1; i < len(ss)-1; i++ {
		text += ss[i] + "\n"
	}
	defer f.Close()
	ioutil.WriteFile(file, []byte(text), 0666)
	fmt.Print("\n\n删除该地址 => " + ss[0])
}
