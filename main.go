package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "shift_corder/routers"
	"strings"

	"time"

	"github.com/astaxie/beego"
)

func main() {
	fmt.Println("自动化办公系统之考勤助手已启动~")
	httpGetMorningRecordOnServer()
	// fmt.Println(beego.AppConfig.String("username"))
	// fmt.Println(beego.AppConfig.String("password"))
	// fmt.Println(beego.AppConfig.Strings("weekend"))
	// a := beego.AppConfig.Strings("weekend")
	// fmt.Println(len(a))
	// for i, v := range a {
	// 	fmt.Println("-------")
	// 	fmt.Printf("i = %d , v = %s\n", i, v)
	// 	fmt.Println("-------")
	// }
	fmt.Println(time.Now())
	//fmt.Println(time.Now().Minute())
	timeInterval, _ := beego.AppConfig.Int64("interval")
	for range time.Tick(time.Duration(timeInterval) * time.Second) {
		fmt.Println("监测中。。。")
		loopThings()
	}

}

func loopThings() {
	t := time.Now()
	if isTodayWeekend(t.String()) {
		fmt.Println("周末一定要睡到中午！")
	} else {
		if checkMorning(t.Hour(), t.Minute()) {
			urlget := beego.AppConfig.String("urlmain")
			// urlget := urlmain
			bodyG, cookies := httpDoGet(urlget)
			data := getInputValues(string(bodyG), true)
			urls := beego.AppConfig.String("urllogin")
			// urls := urllogin
			_, bodyLogin := httpPostNoHeader(cookies, urls, data)
			if checkMorningShiftCorder(string(bodyLogin)) {
				urlattendance := beego.AppConfig.String("urlaction")
				// urlattendance := urlaction
				valuesRecord := getInputValues(string(bodyLogin), false)
				httpPostNoHeader(cookies, urlattendance, valuesRecord)
				fmt.Println("早上签到")
				httpGetMorningRecordOnServer()
			} else {
				fmt.Println("早上签过啦")
			}
		} else if checkAfternoon(t.Hour(), t.Minute()) {
			urlget := beego.AppConfig.String("urlmain")
			// urlget := urlmain

			fmt.Printf("urlget=%s", urlget)

			bodyG, cookies := httpDoGet(urlget)
			data := getInputValues(string(bodyG), true)

			//fmt.Println(data)

			urls := beego.AppConfig.String("urllogin")
			// urls := urllogin

			//fmt.Printf("urls=%s", urls)

			_, bodyLogin := httpPostNoHeader(cookies, urls, data)
			if checkAfternoonShiftCorder(string(bodyLogin)) {
				urlattendance := beego.AppConfig.String("urlaction")
				// urlattendance := urlaction
				valuesRecord := getInputValues(string(bodyLogin), false)
				httpPostNoHeader(cookies, urlattendance, valuesRecord)
				fmt.Println("下午签出")
				httpGetAfternoonRecordOnServer()
			} else {
				fmt.Println("晚上签过啦")
			}

		} else {
			fmt.Println("没到点那~")
		}
	}
}

// func doSomethingOmoshiroi()  {

// 	urlget := "http://kq.neusoft.com"
// 	bodyG, cookies := httpDoGet(urlget)
// 	data := getInputValues(string(bodyG), true)
// 	urls := "http://kq.neusoft.com/login.jsp"
// 	_, bodyLogin := httpPostNoHeader(cookies, urls, data)
// 	urlattendance := "http://kq.neusoft.com/record.jsp"
// 	valuesRecord := getInputValues(string(bodyLogin), false)
// 	httpPostNoHeader(cookies, urlattendance, valuesRecord)

// }
func httpDoGet(url string) ([]byte, []*http.Cookie) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", cookie)
	// req.Header.Set("Host", "kq.neusoft.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:53.0) Gecko/20100101 Firefox/53.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	// req.Header.Set("Referer", "http://kq.neusoft.com/index.jsp?error=3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Pragma", "no-cache")
	// req.Header.Set("Cache-Control", "no-cache")
	respons, err := client.Do(req)
	defer respons.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(respons.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println("GET----BODY---------------------------")
	// fmt.Println(string(body))
	// fmt.Println("GET----HEADER---------------------------")
	// fmt.Println(respons.Header)
	// fmt.Println("GET-----REQUEST---------------------------")
	// fmt.Println(respons.Request)

	// cookieString := respons.Header["Set-Cookie"][0]
	cookies := respons.Cookies()
	return body, cookies
}

func httpDoPost(cookies []*http.Cookie, urls string, paras string) {
	//fmt.Println("POST COOKIES---------------------------------")
	//fmt.Println(cookies)
	client := &http.Client{}
	req, err := http.NewRequest("POST", urls, strings.NewReader(paras))
	//fmt.Println("httpDo---paras---BEGIN------------------------")
	//fmt.Println(strings.NewReader(paras))
	//fmt.Println("httpDo---paras---END--------------------------")
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range cookies {
		req.AddCookie(v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("JSESSIONID", cookie)
	// req.Header.Set("Host", "kq.neusoft.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:53.0) Gecko/20100101 Firefox/53.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Referer", "http://kq.neusoft.com/login.jsp")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Pragma", "no-cache")
	// req.Header.Set("Cache-Control", "no-cache")
	// respons, err := http.PostForm(urls, paras)
	respons, err := client.Do(req)
	defer respons.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = ioutil.ReadAll(respons.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("POST---BODY---------------------------")
	//fmt.Println(string(body))
	// fmt.Println("POST---HEADER---------------------------")
	// fmt.Println(respons.Header)
	// fmt.Println("POST----REQUEST---------------------------")
	// fmt.Println(respons.Request)
}

func httpPostNoHeader(cookies []*http.Cookie, urls string, paras url.Values) ([]*http.Cookie, []byte) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", urls, strings.NewReader(paras.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//fmt.Println("httpDo---paras---BEGIN------------------------")
	//fmt.Println(strings.NewReader(paras.Encode()))
	//fmt.Println("httpDo---paras---END--------------------------")
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range cookies {
		req.AddCookie(v)
		//fmt.Printf("-URL---%s---cookie --%s-----------------------\n", urls, v)
	}
	// uurl := url.URL{}
	// uurl.Path = "login.asp"
	// uurl.Host = "/"
	// uurl.Scheme = "http"
	// fmt.Printf("COOKIE-%s----------------------------------", cookies[0])
	// fmt.Println(uurl)
	// fmt.Println(cookies)
	// client.Jar.SetCookies(&uurl, cookies)
	respons, err := client.Do(req)
	// respons, err := client.PostForm(urls, paras)
	defer respons.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	body, err := ioutil.ReadAll(respons.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println("POST---BODY---------------------------")
	// fmt.Println(string(body))
	// fmt.Println("POST---HEADER---------------------------")
	// fmt.Println(respons.Header)
	// fmt.Println("POST----REQUEST---------------------------")
	// fmt.Println(respons.Request)

	return respons.Cookies(), body
}

// <form action="/login.jsp" method="post" name="LoginForm">
// 				<input type="hidden" name="login" value="true" />
func getNameAndValue(src string) (string, string) {
	tmp := strings.Index(src, `"`)
	name := src[:tmp]
	var value = ""
	if strings.Contains(src, "value=") {
		tmpArray := strings.Split(src, `value="`)
		vtmp := strings.Index(tmpArray[1], `"`)
		value = tmpArray[1][:vtmp]
	}
	return name, value
}

func getInputValues(src string, needUser bool) url.Values {
	inputArray := strings.Split(src, `name="`)
	// fmt.Println("getInputValues---------------------------")
	// fmt.Println(inputArray)
	// fmt.Println("getInputValues---------------------------")
	var data = make(url.Values)

	for i, v := range inputArray {
		if i <= 1 {
			continue
		}
		name, value := getNameAndValue(v)
		// fmt.Println("LOOP-BEGIN---------------------------")
		// fmt.Printf("index = %d, name = %s, value = %s\n", i, name, value)
		// fmt.Println("LOOP--END----------------------------")
		if needUser {
			if i == len(inputArray)-2 {
				data[name] = []string{beego.AppConfig.String("username")}
				// data[name] = []string{username}
				continue
			}
			if i == len(inputArray)-1 {
				data[name] = []string{beego.AppConfig.String("password")}
				// data[name] = []string{password}
				continue
			}
		}
		data[name] = []string{value}

	}
	// fmt.Println("Input----Values-BEGIN-----------------------------")
	// fmt.Println(data)
	// fmt.Println("Input----Values-END-------------------------------")
	return data
}

//可以打就是true
func checkMorningShiftCorder(src string) bool {
	srcArray := strings.Split(src, "tbody")
	//目测观察tbody只有一对，所以分出的数组一定是三个字符串，只要中间的。
	dst := srcArray[1]
	//判断里面有几个tr，早上就要一对就行了,没有就true,多了不管~
	if strings.Contains(dst, "tr") {
		if strings.Count(dst, "td") > 3 {
			return false
		}
	}
	return true

}

//可以打就是true
func checkAfternoonShiftCorder(src string) bool {
	srcArray := strings.Split(src, "tbody")
	//目测观察tbody只有一对，所以分出的数组一定是三个字符串，只要中间的。
	dst := srcArray[1]
	//判断里面有几个tr，晚上没有或一个就true就行了,多了不管~
	if strings.Count(dst, "tr") > 2 {
		return false
	}
	return true
}

func isTodayWeekend(src string) bool {
	for _, v := range beego.AppConfig.Strings("weekend") {
		// for _, v := range weekend {
		//fmt.Println(v)
		if strings.ToLower(src) == v {
			return true
		}
	}
	return false
}

func checkMorning(hour int, minute int) bool {
	confHour, err := beego.AppConfig.Int("morninghour")
	if err != nil {
		fmt.Println("Morning hour setting format error!")
	}
	confMinute, err := beego.AppConfig.Int("morningminute")
	if err != nil {
		fmt.Println("Morning minute setting format error!")
	}
	if hour == confHour && minute >= confMinute {
		return true
	}
	return false
}

func checkAfternoon(hour int, minute int) bool {

	confHour, err := beego.AppConfig.Int("afternoonhour")
	if err != nil {
		fmt.Println("Afternoon hour setting format error!")
	}
	confMinute, err := beego.AppConfig.Int("afternoonminute")
	if err != nil {
		fmt.Println("Afternoon minute setting format error!")
	}

	if hour == confHour && minute >= confMinute {
		return true
	}
	return false
}
//早上打卡成功通知到俺的服务器。
func httpGetMorningRecordOnServer() {
	urlproxy := beego.AppConfig.String("urlproxy")
	url1, _ := url.Parse(urlproxy)
	proxy := http.ProxyURL(url1)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}
	//client := &http.Client{}
	url := fmt.Sprintf("http://47.xx.xx.38:8777/api/1.0/morningshiftcoder?user=%s", beego.AppConfig.String("username"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.Do(req)
}
//下午打卡通知到俺的服务器
func httpGetAfternoonRecordOnServer() {
	urlproxy := beego.AppConfig.String("urlproxy")
	url1, _ := url.Parse(urlproxy)
	proxy := http.ProxyURL(url1)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}
	//client := &http.Client{}
	url := fmt.Sprintf("http://47.100.xx.xx:8777/api/1.0/afternoonshiftcoder?user=%s", beego.AppConfig.String("username"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.Do(req)
}
