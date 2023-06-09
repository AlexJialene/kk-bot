package service

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func Test1(t *testing.T) {
	StartNewsTickService(func(str string) {
		fmt.Println(str)
	})

}

func Test2(t *testing.T) {
	weather := Weather("广州")
	fmt.Println(weather)

}

func Test3(t *testing.T) {
	//StartPicDayService(func(file *os.File) {
	//	fmt.Println("======")
	//	fmt.Println(file.Name())
	//})

}

func Test4(t *testing.T) {
	StartMoyuPicDayService(func(name string) {
		fmt.Println(name)
	})

}

type DemoCall struct {
}

func (d DemoCall) Call(r string) bool {
	fmt.Println("rev callback")
	fmt.Println(r)
	return true
}

func Test5(t *testing.T) {
	CreateCLSRoll(DemoCall{})
	//test
	select {}

}

func Test6(t *testing.T) {
	url := "https://www.cls.cn/nodeapi/refreshTelegraphList?app=CailianpressWeb&lastTime=1686280109&os=web&sv=7.7.5"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Host", "www.cls.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.41")
	request.Header.Add("Referer", "https://www.cls.cn/telegraph")
	request.Header.Add("Cookie", "HWWAFSESID=6a957d9e8ad30becbd; HWWAFSESTIME=1686280146130; hasTelegraphNotification=on; hasTelegraphRemind=on; hasTelegraphSound=on; vipNotificationState=on;")
	//request.AddCookie(&http.Cookie{Name: "", Value: ""})
	client := http.Client{}
	get, _ := client.Do(request)
	//get, _ := http.Get("https://www.cls.cn/v1/roll/get_roll_list?app=CailianpressWeb&category=red&last_time=1686277039&os=web&refresh_type=1&rn=20&sv=7.7.5")
	all, _ := io.ReadAll(get.Body)
	defer get.Body.Close()
	fmt.Println(all)
	fmt.Println(string(all))
}

func Test7(t *testing.T) {
	url := "https://www.cls.cn/v1/roll/get_roll_list?app=CailianpressWeb&category=red&last_time=1686277039&os=web&refresh_type=1&rn=20&sv=7.7.5"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Host", "www.cls.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.41")
	request.Header.Add("Referer", "https://www.cls.cn/telegraph")
	request.Header.Add("Cookie", "HWWAFSESID=6a957d9e8ad30becbd; HWWAFSESTIME=1686280146130; hasTelegraphNotification=on; hasTelegraphRemind=on; hasTelegraphSound=on; vipNotificationState=on;")
	//request.AddCookie(&http.Cookie{Name: "", Value: ""})
	client := http.Client{}
	get, _ := client.Do(request)
	//get, _ := http.Get("https://www.cls.cn/v1/roll/get_roll_list?app=CailianpressWeb&category=red&last_time=1686277039&os=web&refresh_type=1&rn=20&sv=7.7.5")
	all, _ := io.ReadAll(get.Body)
	defer get.Body.Close()
	fmt.Println(all)
	fmt.Println(string(all))
}
