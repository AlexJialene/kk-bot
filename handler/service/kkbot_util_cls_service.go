package service

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Callback interface {
	Call(r string)
}

type ResponseBody struct {
	Error int          `json:"error"`
	Data  ResponseData `data`
}

type ResponseData struct {
	Data []RollData `roll_data`
}

type RollData struct {
	Content string `json:"content"`
	Id      int    `json:"id"`
	Level   string `json:"level"`
}

var funcList []Callback
var lastCLSId int

func CreateCLSRoll(callback Callback) {
	funcList = append(funcList, callback)

}

func init() {
	lastCLSId = 0
	funcList = []Callback{}
	log.Println("start to load the roll_data")
	//go func() {
	//
	//	for {
	//		fmt.Println(111)
	//		//test
	//		time.Sleep(10 * time.Second)
	//
	//		result := GetCLSRollList()
	//		assemble(result)
	//		time.Sleep(10 * time.Minute)
	//	}
	//}()
}

func assemble(result *ResponseBody) {
	if result == nil {
		return
	}
	if result.Error == 0 {
		if len(result.Data.Data) > 0 {
			for _, v := range result.Data.Data {
				if len(v.Content) > 0 && len(v.Level) > 0 {
					if v.Id > lastCLSId {
						lastCLSId = v.Id
						for _, f := range funcList {
							f.Call(v.Content)
						}
					}
				}
			}
		}
	}
}

func GetCLSRollList() *ResponseBody {

	//todo 2023/6/9 lamkeizyi -

	client := http.Client{}
	//client.

	unix := time.Now().Unix()
	sprintf := fmt.Sprintf("app=CailianpressWeb&category=red&last_time=%d&os=web&refresh_type=1&rn=20&sv=7.7.5", unix)
	sha1Hash := sha1.Sum([]byte(sprintf))
	sum := fmt.Sprintf("%x", sha1Hash)

	bytes := md5.Sum([]byte(sum))
	md := fmt.Sprintf("%x", bytes)
	fmt.Println("sign = ", md)

	url := fmt.Sprintf("https://www.cls.cn/v1/roll/get_roll_list?app=CailianpressWeb&category=red&last_time=%d&os=web&refresh_type=1&rn=20&sv=7.7.5&sign=%s", unix, md)
	//url := "" + md
	//trueUrl := fmt.Sprintf(url, time.Now().Unix())
	request, _ := http.NewRequest("GET", url, nil)
	//request.AddCookie(&http.Cookie{Name: "HWWAFSESID", Value: "2c114bb22a69f4e145"})
	//request.AddCookie(&http.Cookie{Name: "HWWAFSESTIME", Value: "1686236077500"})
	request.Header.Add("Host", "www.cls.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Accept-Encoding", "gzip, deflate, br, json")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Referer", "https://www.cls.cn/telegraph")

	get, err := client.Do(request)
	//get, err := http.Get(trueUrl)
	if err != nil {
		log.Println("get cls rollList has error ", err)
	}
	defer get.Body.Close()
	all, err := io.ReadAll(get.Body)
	fmt.Println(url)
	fmt.Println(all)
	fmt.Println(1123123)
	s := string(all)
	log.Printf("all = %s", s)
	fmt.Println(string(all))
	if err != nil {
		log.Printf("get cls rollList has error, body = %s \n ", string(all))
	}
	r := &ResponseBody{}
	if err = json.Unmarshal(all, r); err == nil {
		return r
	}

	return nil
}
