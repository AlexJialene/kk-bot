package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Callback callback 需要获取财联社电报需要实现此接口
type Callback interface {
	Call(r string) bool
}

type DefaultCallbackAdapter struct {
}

func (a DefaultCallbackAdapter) Call(r string) bool {
	log.Println("not doing anything ")
	return true
}

type ResponseBody struct {
	Error int          `json:"error"`
	Data  ResponseData `json:"data"`
}

type ResponseData struct {
	Data []RollData `json:"roll_data"`
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
	go func() {

		for {
			//test
			fmt.Println(111)
			time.Sleep(10 * time.Second)

			result := GetCLSRollList()
			assemble(result)

			fmt.Println("ending ")
			time.Sleep(10 * time.Minute)
		}
	}()
}

func assemble(result *ResponseBody) {
	if result == nil {
		return
	}

	if result.Error == 0 {
		if len(result.Data.Data) > 0 {
			for i := len(result.Data.Data) - 1; i >= 0; i-- {
				v := result.Data.Data[i]

				//fmt.Printf("%d  id  = %d", i, v.Id)
				//fmt.Printf("%d = content = %s", i, v.Content[0:10])
				//fmt.Println("======")
				//fmt.Printf("id = %d , level = %s", v.Id, v.Level)
				//fmt.Println("======")

				if v.Id != 0 && len(v.Level) > 0 {
					//fmt.Println("------")
					//fmt.Printf("id = %d , lastCLSId = %d", v.Id, lastCLSId)
					//fmt.Println("------")
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
	url := "https://www.cls.cn/v1/roll/get_roll_list?app=CailianpressWeb&category=red&last_time=%d&os=web&refresh_type=1&rn=20&sv=7.7.5"
	sprintf := fmt.Sprintf(url, time.Now().Unix())
	get := clsGet(sprintf)

	if get != nil && len(get) > 0 {
		r := &ResponseBody{}
		err := json.Unmarshal(get, r)
		if err != nil {
			log.Println("parse json error , ", err)
			return nil
		}

		//fmt.Println("-------")
		//fmt.Println(r)
		//fmt.Println(r.Data)
		//fmt.Println(r.Data.Data)
		//for _, v := range r.Data.Data {
		//	fmt.Println(v.Id)
		//	fmt.Println(v.Level)
		//}
		//fmt.Println("-------")
		return r
	}
	return nil
}

func clsGet(url string) []byte {
	client := http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Host", "www.cls.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.41")
	request.Header.Add("Referer", "https://www.cls.cn/telegraph")
	request.Header.Add("Cookie", "HWWAFSESID=6a957d9e8ad30becbd; HWWAFSESTIME=1686280146130; hasTelegraphNotification=on; hasTelegraphRemind=on; hasTelegraphSound=on; vipNotificationState=on;")
	do, err := client.Do(request)
	if err != nil {
		log.Println("request www.cls.cn error , ", err)
		return nil
	}
	defer do.Body.Close()
	all, err := io.ReadAll(do.Body)
	if err != nil {
		log.Println("request www.cls.cn get body error , ", err)
		return nil
	}
	return all
}
