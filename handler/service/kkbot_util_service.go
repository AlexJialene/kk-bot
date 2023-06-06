package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	loader "sf-bot/handler/load"
	"strconv"
	"strings"
	"time"
)

type BaseVO struct {
	Code int32
	Data []interface{}
}

type WeatherVO struct {
	City    string
	Weather string
	Date    string
}

type ZhihuVO struct {
	Name string
	Hot  string
	Url  string
}

type EveryDayText struct {
	ImgUrl string
	Name   string
	Time   []string
	Data   []string
}

func (e *EveryDayText) ToString() string {
	r := ""
	r = r + e.Name + "\n"
	for _, v := range e.Time {
		r = r + v + " "
	}
	r = r + "\n \n"
	for i, datum := range e.Data {
		r = r + strconv.Itoa(i+1) + "、" + datum + "\n \n"
	}
	return r
}

type Moyu struct {
	Success bool
	Url     string
}

// StartNewsTickService 知乎热版
func StartNewsTickService(f func(str string)) {
	//t := "0 0 9 * * ?"
	//t := "0 0 8 ? * 2,3,4,5,6 *"
	//c := cron.New()
	//c.AddFunc(t, func() {
	client := http.Client{Timeout: 60 * time.Second}
	if get, err := client.Get("https://tenapi.cn/v2/zhihuhot"); err == nil {
		defer get.Body.Close()
		all, _ := io.ReadAll(get.Body)
		z := &BaseVO{}
		json.Unmarshal(all, z)
		if z.Code == 200 {
			result := ""
			//log.Printf("the zhihu hot list = %s \n", z.Data)
			for i, v := range z.Data {
				vo := &ZhihuVO{}
				interfaceToStruct(v.(map[string]interface{}), vo)
				result = result + "" + strconv.Itoa(i+1) + "." + strings.Trim(vo.Name, " ") + "\n \n"
			}
			f(result)
		}
	}

	//})

}

func StartPicTickService(f func(file *os.File)) {

}

func GetPicDayTextService() (*EveryDayText, error) {
	client := http.Client{Timeout: 60 * time.Second}
	form, _ := client.Get("https://api.vvhan.com/api/60s?type=json")
	if all, err := io.ReadAll(form.Body); err == nil {
		defer form.Body.Close()
		b := &EveryDayText{}
		err := json.Unmarshal(all, b)
		return b, err
	}
	return nil, errors.New("GetPicDayTextService error ")
}

func StartPicDayService(f func(s string)) {
	if service, err := GetPicDayTextService(); err == nil {
		if len(service.ImgUrl) > 0 {
			client := http.Client{Timeout: 60 * time.Second}
			if get, err := client.Get(service.ImgUrl); err == nil {
				readAll, _ := io.ReadAll(get.Body)
				fileName := openFileName()
				if err := os.WriteFile(fileName, readAll, 0777); err == nil {
					f(fileName)
				}
				defer get.Body.Close()
			}
		}
	}

}

func StartMoyuPicDayService(f func(name string)) {
	log.Println("load moyu pic ... ")
	client := http.Client{Timeout: 60 * time.Second}
	form, _ := client.Get("https://api.vvhan.com/api/moyu?type=json")
	if all, err := io.ReadAll(form.Body); err == nil {
		defer form.Body.Close()
		b := &Moyu{}
		json.Unmarshal(all, b)
		if len(b.Url) > 0 {
			fmt.Printf("url  = %s \n ", b.Url)
			//b.Url = "https://img1.baidu.com/it/u=2932772446,2464263128&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=1082"
			if get, err := http.Get(b.Url); err == nil {
				readAll, _ := io.ReadAll(get.Body)
				name := openFileName()
				if err := os.WriteFile(openFileName(), readAll, 0777); err == nil {
					log.Println("callback pic ... ")
					f(name)
				}
				defer get.Body.Close()
			}
		}
	}
}

// Weather 天气
func Weather(city string) *WeatherVO {
	client := http.Client{Timeout: 60 * time.Second}
	u := url.Values{}
	u.Add("city", city)
	form, _ := client.PostForm("https://tenapi.cn/v2/weather", u)
	if all, err := io.ReadAll(form.Body); err == nil {
		defer form.Body.Close()
		b := &BaseVO{}
		json.Unmarshal(all, b)
		if b.Code == 200 {
			vo := &WeatherVO{}
			interfaceToStruct(b.Data[0].(map[string]interface{}), vo)
			return vo
		}
	}
	return nil
}

func openFileName() string {
	load := loader.Load("commons.file_temp_dir")
	goos := runtime.GOOS
	sep := ""
	s := strconv.FormatInt(time.Now().Unix(), 10)
	if goos == "windows" {
		sep = "\\"
	}
	if goos == "linux" {
		sep = "/"
	}

	if len(load) > 0 {
		log.Printf("the temp file dir = %s \n ", load)
		return load + sep + s + ".png"
	} else {
		log.Printf("the temp file dir = %s \n ", os.TempDir())
		return os.TempDir() + sep + s + ".png"
	}
}

func interfaceToStruct(m map[string]interface{}, obj interface{}) error {
	for k, v := range m {
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	trueName := strings.ToUpper(name[:1]) + name[1:]
	structFieldValue := structValue.FieldByName(trueName)
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}
	structFieldValue.Set(val)
	return nil
}
