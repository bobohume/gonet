package base

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type(
	Kv struct {
		Key string
		Val interface{}
	}
)

//不能用map,map是桶,无序,有些get请求需要时序
func geturl(url string, params ...Kv) string{
	str := ""
	for i, v := range params{
		if i > 0{
			str += "&"
		}
		switch  typ := v.Val.(type) {
		case int8:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(int8))
		case int16:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(int16))
		case int32:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(int32))
		case int:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(int))
		case int64:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(int64))
		case uint8:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(uint8))
		case uint16:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(uint16))
		case uint32:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(uint32))
		case uint:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(uint))
		case uint64:
			str += fmt.Sprintf("%s=%d", v.Key, v.Val.(uint64))
		case float32:
			str += fmt.Sprintf("%s=%f", v.Key, v.Val.(float32))
		case float64:
			str += fmt.Sprintf("%s=%f", v.Key, v.Val.(float64))
		case string:
			str += fmt.Sprintf("%s=%s", v.Key, v.Val.(string))
		default:
			log.Printf("Url not support type [%d]", typ)
		}
	}
	str = url + "?" + str
	return str
}

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string, params ...Kv) ([]byte, error){
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(geturl(url))
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return res, nil
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data []byte, contentType string) ([]byte, error){
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if contentType == ""{
		contentType = "text/plain; charset=utf-8"
	}
	req.Header.Add("content-type", contentType)
	if err != nil {
		return []byte{}, err
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}