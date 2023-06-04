package handler

import (
	"fmt"
	loader "sf-bot/handler/load"
	"testing"
)

func Test1(t *testing.T) {
	load := loader.Load("gpt.host")
	fmt.Println(load)
	initAgent := InitAgent()
	ask := initAgent.AskAgent().Ask("123123", "静夜思")
	fmt.Println(ask)
	select {}
}
