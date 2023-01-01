package main

import (
	"fmt"
	"github.com/yuzhiyizhan/requests"
	"github.com/yuzhiyizhan/requests/models"
	"github.com/yuzhiyizhan/requests/url"
	"sync"
	"time"
)

var wg sync.WaitGroup

func get_request(urls string, result chan *models.Response) {
	defer wg.Done()
	req := url.NewRequest()
	headers := map[string]interface{}{
		"Accept":          "application/json",
		"Accept-Language": "zh-CN,zh-TW;q=0.9,zh;q=0.8,en-US;q=0.7,en;q=0.6",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	}
	cookies := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	req.Cookies = url.ParseDictCookies(urls, cookies)
	req.Headers = url.ParseDictHeaders(headers)
	req.Proxies = "http://127.0.0.1:7890"
	req.Ja3 = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
	response, err := requests.Get(urls, req)
	if err != nil {
		return
	}
	result <- response
}

func main() {
	start := time.Now()
	var number = []string{}
	result := make(chan *models.Response, 100)
	var arr [100]int
	for i := 0; i < 100; i++ {
		arr[i] = i + 1
	}
	for _, i := range arr {
		fmt.Println(i)
		urls := "https://tls.peet.ws/api/all"
		go get_request(urls, result)
		wg.Add(1)
	}
	go func() {
		for {
			response := <-result
			ip := response.GoJson("ip")
			number = append(number, ip)
			fmt.Println(response.GoJson("ip"))

		}
	}()
	wg.Wait()

	defer func() {
		fmt.Println("成功的请求数: ", len(number))
	}()
	defer func() {
		elapsed := time.Since(start)
		fmt.Println("该函数执行完成耗时：", elapsed)
	}()
}
