package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Danny-Dasilva/fhttp"
	"github.com/antchfx/htmlquery"
	"github.com/bitly/go-simplejson"
	"github.com/tidwall/gjson"
	"github.com/yuzhiyizhan/requests/url"
	"io"
	"os"
	"strings"
)

// Response结构体
type Response struct {
	Url         string
	Headers     http.Header
	Cookies     []*http.Cookie
	Text        string
	Content     []byte
	Body        io.ReadCloser
	StatusCode  int
	History     []*Response
	Request     *url.Request
	xpathResult []string
}

// 使用自带库JSON解析
func (res *Response) Json() (map[string]interface{}, error) {
	js := make(map[string]interface{})
	err := json.Unmarshal(res.Content, &js)
	return js, err
}

// 使用go-simplejson解析
func (res *Response) SimpleJson() (*simplejson.Json, error) {
	return simplejson.NewFromReader(res.Body)
}

// 状态码是否错误
func (res *Response) RaiseForStatus() error {
	var err error
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		err = errors.New(fmt.Sprintf("%d Client Error", res.StatusCode))
	} else if res.StatusCode >= 500 && res.StatusCode < 600 {
		err = errors.New(fmt.Sprintf("%d Server Error", res.StatusCode))
	}
	return err
}

// 使用gjson解析
func (res *Response) GoJson(path string) string {
	data := gjson.Get(res.Text, path).String()
	return data
}

// 使用xpath解析
func (r *Response) Xpath(expr string) *Response {
	r.xpathResult = make([]string, 0)
	doc, err := htmlquery.Parse(strings.NewReader(r.Text))
	if err != nil {
		fmt.Println(err.Error())
		return r
	}
	for _, node := range htmlquery.Find(doc, expr) {
		if len(node.Data) > 0 {
			r.xpathResult = append(r.xpathResult, node.Data)
		}
	}
	return r
}

func (r *Response) Get() (string, bool) {
	if len(r.xpathResult) == 0 {
		return "", false
	}
	return r.xpathResult[0], true
}

func (r *Response) Getall() []string {
	return r.xpathResult
}

// 下载文件
func (r *Response) SaveFile(filename string) error {
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(r.Content)
	if err != nil {
		return err
	}
	return nil
}
