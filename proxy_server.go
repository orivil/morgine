package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		URL := request.URL.Query().Get("url")
		if URL != "" {
			logPrintf("url: %s\n\n", URL)
			req, err := http.NewRequest(request.Method, URL, nil)
			if err != nil {
				logPrintln(err)
				http.Error(writer, err.Error(), http.StatusBadRequest)
			} else {
				for key, value := range request.Header {
					req.Header[key] = value
				}
				logPrintf("request headers: %v\n\n", request.Header)
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					logPrintln(err)
					http.Error(writer, err.Error(), http.StatusBadRequest)
				} else {
					writer.WriteHeader(resp.StatusCode)
					for key, value := range resp.Header {
						writer.Header()[key] = value
					}
					logPrintf("response headers: %v\n\n", resp.Header)
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						http.Error(writer, err.Error(), http.StatusBadRequest)
					}
					//Detect if the content is gzipped
					contentType := http.DetectContentType(data)
					//If we detect gzip, then make a gzip reader, then wrap it in a scanner
					if strings.Contains(contentType, "x-gzip") {
						reader, err := gzip.NewReader(bytes.NewBuffer(data))
						if err != nil {
							logPrintln(err)
							http.Error(writer, err.Error(), http.StatusBadRequest)
							return
						}
						data, err = ioutil.ReadAll(reader)
						if err != nil {
							logPrintln(err)
							http.Error(writer, err.Error(), http.StatusBadRequest)
						}
					}

					logPrintf("response code [%d] body: \n%s\n\n", resp.StatusCode, string(data))
					_, _ = writer.Write(data)
				}
			}
		}
	}
	log.Fatal(http.ListenAndServe(":8083", handler))
}

func logPrintf(fmt string, v ...interface{}) {
	open := os.Getenv("open_log")
	if open != "" {
		log.Printf(fmt, v...)
	}
}

func logPrintln(v ...interface{}) {
	open := os.Getenv("open_log")
	if open != "" {
		log.Println(v...)
	}
}
