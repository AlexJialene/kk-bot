package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
	loader "sf-bot/handler/load"
	"strings"
)

// var Customize bool = false
var agent *AgentFactory
var wx *openwechat.Self

const (
	aiteMe      = "@alex_kkbot"
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
	//dispatcher.OnGroup(func(ctx *openwechat.MessageContext) {
	//	fmt.Printf("onGroup message = %s", ctx.Content)
	//})

	dispatcher.OnUser(func(user *openwechat.User) bool {
		fmt.Printf("the on user = %s", user)
		return user.IsGroup() && loader.Exist(user.NickName)

	}, func(ctx *openwechat.MessageContext) {
		log.Printf("receive the msg = %s \n", ctx.Content)

		if ctx.IsText() {
			if strings.Contains(ctx.Content, aiteMe) {
				user, _ := ctx.SenderInGroup()

				fmt.Printf("the user_info = %s \n", user)

				id := user.NickName
				fmt.Printf("nickName = %s \n ", id)
				//fmt.Printf("userName = %s \n ", user.UserName)
				//marshal, _ := json.Marshal(user)
				//fmt.Printf("json = %s", string(marshal))

				msg := strings.ReplaceAll(ctx.Content, aiteMe, "")
				answer := agent.AskAgent().Ask(groupPrefix+id, msg)
				_, err := ctx.ReplyText("@" + id + " \n" + answer)
				if err != nil {
					log.Println(err)
					return
				}
			}

		} else {
			//todo 2023/6/4 lamkeizyi - 其他消息暂未支持

		}

	})

	bot.MessageHandler = dispatcher.AsMessageHandler()

}
