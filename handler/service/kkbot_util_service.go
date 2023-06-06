package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type EveryDay struct {
	ImgUrl string
}

type Moyu struct {
	success bool
	url     string
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

func StartPicDayService(f func(file *os.File)) {
	client := http.Client{Timeout: 60 * time.Second}
	form, _ := client.Get("https://api.vvhan.com/api/60s?type=json")
	if all, err := io.ReadAll(form.Body); err == nil {
		defer form.Body.Close()
		b := &EveryDay{}
		json.Unmarshal(all, b)
		if len(b.ImgUrl) > 0 {
			open, err := os.Create(openFileName())
			if err != nil {
				fmt.Println(err)
			} else {
				if get, err := client.Get(b.ImgUrl); err == nil {
					readAll, _ := io.ReadAll(get.Body)
					if _, err := open.Write(readAll); err == nil {
						f(open)
					}
					defer get.Body.Close()
					defer open.Close()
				}
			}
		}
	}
}

func StartMoyuPicDayService(f func(file *os.File)) {
	client := http.Client{Timeout: 60 * time.Second}
	form, _ := client.Get("https://api.vvhan.com/api/moyu?type=json")
	if all, err := io.ReadAll(form.Body); err == nil {
		defer form.Body.Close()
		b := &Moyu{}
		json.Unmarshal(all, b)
		if len(b.url) > 0 {
			open, err := os.Create(openFileName())
			if err != nil {
				fmt.Println(err)
			} else {
				if get, err := client.Get(b.url); err == nil {
					readAll, _ := io.ReadAll(get.Body)
					if _, err := open.Write(readAll); err == nil {
						f(open)
					}
					defer get.Body.Close()
					defer open.Close()
				}
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
	s := string(time.Now().Unix())
	if goos == "windows" {
		sep = "\\"
	}
	if goos == "linux" {
		sep = "/"
	}
	if len(load) > 0 {
		return load + sep + s + ".png"
	} else {
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
