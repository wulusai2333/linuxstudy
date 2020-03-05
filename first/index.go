package main

import (
	"fmt"
	"github.com/wulusai2333/gostudy/logger"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var requestLog = make(chan *string, 10000)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	err := http.ListenAndServeTLS(":443", "1_wulusai.net_bundle.crt", "2_wulusai.net.key", mux)
	if err != nil {
		fmt.Println("http server start failed,err:", err)
	}
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", index)
		err := http.ListenAndServe(":80", mux)
		if err != nil {
			fmt.Println("http server start failed,err:", err)
		}
	}()
	go func() {
		for {
			select {
			case logStr, ok := <-requestLog:
				if !ok {
					return
				}
				logs := logger.NewFileLogger(logger.INFO, "/root/github.com/linuxstudy/first/", "wulusai", 1024*1024*10)
				logs.Info(*logStr)
			default:
				time.Sleep(time.Second)
			}
		}
	}()
}

func index(writer http.ResponseWriter, request *http.Request) {
	file, err := os.Open("./index.html")
	if err != nil {
		fmt.Println("open file failed,err:", err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file failed,err:", err)
		return
	}
	i, err := writer.Write(data)
	if err != nil {
		fmt.Println("response write failed,err:", err)
		return
	}
	requestString := fmt.Sprintln(request.Host, request.Method, request.Header, i)
	requestLog <- &requestString
}
