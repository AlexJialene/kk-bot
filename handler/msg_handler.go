package handler

import (
	"github.com/eatmoreapple/openwechat"
	"log"
)

// var Customize bool = false
var agent *AgentFactory
var wx *openwechat.Self

const (
	//aiteMe      = "@alex_kkbot"
	groupPrefix = "gr"
)

func Bootstrap(bot *openwechat.Bot, weixin *openwechat.Self) {
	wx = weixin
	agent = InitAgent()
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	groups, err := weixin.Groups()
	if err != nil {
		log.Panic(err)
		return
	}

	log.Println(groups)
	dispatcher := openwechat.NewMessageMatchDispatcher()
	bootStrapDispatcher(dispatcher)
	bot.MessageHandler = dispatcher.AsMessageHandler()

}

func bootStrapDispatcher(dispatcher *openwechat.MessageMatchDispatcher) {
	// 2023/6/8 lamkeizyi - group inits
	groupHandler := CreateGroupBotHandler()
	if groupHandler != nil {
		dispatcher.OnUser(func(user *openwechat.User) bool { return user.IsGroup() && groupHandler.exists(user.NickName) }, groupHandler.recv)
	}
	// 2023/6/8 lamkeizyi - user inits
	uh := CreateUserHandler()
	if uh != nil {
		dispatcher.OnUser(func(user *openwechat.User) bool { return user.IsFriend() }, uh.Recv)
	}

}
