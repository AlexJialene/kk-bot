package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
	loader "sf-bot/handler/load"
	"strings"
)

var groupHandler *GroupBotHandler

type GroupBotHandler struct {
	closeReplySuffix bool
	groupNames       []string
}

// return is break
func (g *GroupBotHandler) recvCommand(command string, f func(i string)) bool {
	if strings.Contains(command, "#关闭后缀#") {
		g.closeReplySuffix = true
		f(g.infos())
		return true
	}

	if strings.Contains(command, "#开启后缀#") {
		g.closeReplySuffix = false
		f(g.infos())
		return true
	}

	if !g.closeReplySuffix {
		f(command + "\n" + "------------\n" + "已收到，我知道你很急，但是你先别急！有问题联系Alex_")
	}
	return false
}

func (g *GroupBotHandler) exists(name string) bool {
	split := strings.Split(loader.GroupName(), ",")
	for _, key := range split {
		if strings.Contains(name, key) {
			return true
		}
	}
	return false
}

func (g *GroupBotHandler) infos() string {
	info := ""
	info = info + "----setting----\n后缀："
	if g.closeReplySuffix {
		info += "已关闭\n"
	} else {
		info += "已开启\n"
	}
	info = info + "----setting----"
	return info
}

func (g *GroupBotHandler) recv(ctx *openwechat.MessageContext) {
	if ctx.IsText() {
		log.Printf("receive the msg = %s \n", ctx.Content)

		if strings.Contains(ctx.Content, aiteMe) {
			user, _ := ctx.SenderInGroup()
			fmt.Printf("the user_info = %s \n", user)
			id := user.NickName
			fmt.Printf("nickName = %s \n ", id)

			msg := strings.ReplaceAll(ctx.Content, aiteMe, "")
			command := g.recvCommand(msg, func(info string) { ctx.ReplyText(info) })
			if command {
				return
			}
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
}

func CreateGroupBotHandler() *GroupBotHandler {
	groupHandler = &GroupBotHandler{closeReplySuffix: false, groupNames: strings.Split(loader.GroupName(), ",")}
	return groupHandler
}
