package gpt

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	_ = newSender("id1")
	fmt.Println(1)
	time.Sleep(20 * time.Second)
	fmt.Println(2)
	_ = newSender("id2")
	select {}
}

func Test2(t *testing.T) {
	u := uuid.New()
	fmt.Println(u.String())
	all := strings.ReplaceAll(u.String(), "-", "")
	fmt.Println(all)
}

func Test3(t *testing.T) {
	//host := "http://127.0.0.1:19090"
	//gptHost := CreateGptHost(host)
	//s := gptHost.Talk("alex_test_id_4", "静夜思全文")
	//fmt.Println(s)
	//select {}
}

func Test4(t *testing.T) {
	//host := "http://127.0.0.1:19090"
	//CreateGptHost(host)
	//title := genTitle("a1fc8d52-d653-465d-ab26-b562fb485aed", "gpt-3.5-turbo", "60a5ca1e-753a-4e81-9971-f6c346dfbdb8")
	//fmt.Println(title)
}
