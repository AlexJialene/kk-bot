package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"sf-bot/handler"
)

func main() {
	boot()
}

func boot() {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	handler.Bootstrap(bot, self)
	err = bot.Block()
	if err != nil {
		panic(err)
	}
}
