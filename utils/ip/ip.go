// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package ip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// 获得域名 IP
func GetWebSiteIP(domain string) (ip string, err error) {
	var ls []string
	ls, err = net.LookupHost(domain)
	if err != nil {
		return "", err
	} else {
		for _, ip = range ls {
			return ip, nil
		}
	}
	return "", fmt.Errorf("domain %s get ip field", domain)
}

type sohuCitySN struct {
	CIP   string `json:"cip"`
	CID   string `json:"cid"`
	CName string `json:"cname"`
}

func GetRemoteIP() (string, error) {
	resp, err := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	jsonStart := bytes.Index(content, []byte("{"))
	jsonEnd := bytes.Index(content, []byte("}"))
	if 0 <= jsonStart && jsonStart < jsonEnd {
		content = content[jsonStart : jsonEnd+1]
		csn := &sohuCitySN{}
		err = json.Unmarshal(content, csn)
		if err != nil {
			return "", err
		}
		return csn.CIP, nil
	} else {
		return "", fmt.Errorf("got content: %s", content)
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("unexpected error")
}

// TODO: 待验证
func GetHttpRequestIP(request *http.Request) string {
	addr := request.Header.Get("X-Real-IP")
	if addr == "" {
		addr = request.Header.Get("X-Forwarded-For")
		if addr != "" {
			addrs := strings.Split(addr, ",")
			for _, value := range addrs {
				if len(value) > 0 {
					addr = value
					break
				}
			}
		}
	}
	if addr == "" {
		addr = request.RemoteAddr
	}
	return addr
}
