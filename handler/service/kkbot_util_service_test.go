package service

import (
	"fmt"
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

func (d DemoCall) Call(r string) {
	fmt.Println(r)
}

func Test5(t *testing.T) {
	//CreateCLSRoll(DemoCall{})
	//select {}
	//list := GetCLSRollList()
	//fmt.Println(list)

}
