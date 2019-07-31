// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	//proxy := "http://127.0.0.1:8083"
	proxy := ""
	//req, err := http.NewRequest("GET", "https://wx487e31190600d31e.zhixuaw.cn/index/book/info?book_id=11000162465", nil)
	//req, err := http.NewRequest("GET", initUrl(proxy, "http://www.baidu.com?id=123"), nil)
	req, err := http.NewRequest("GET", initUrl(proxy, "https://wx487e31190600d31e.zhixuaw.cn/index/book/info?book_id=11000006458"), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat QBCore/3.43.884.400 QQBrowser/9.0.2524.400")
	req.Header.Set("Cookie", "sex=girl; back=; user_id=13911255; channel_id=8923; openid=oCq6ivwcu4bQ2_QxthXwCQtkaUho;")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err, resp)
	}
	fmt.Println(resp.Header)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func initUrl(proxy string, Url string) string {
	if proxy != "" {
		return proxy + "/proxy?url=" + url.QueryEscape(Url)
	} else {
		return Url
	}
}
