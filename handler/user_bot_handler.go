package handler

import (
	"encoding/json"
	"github.com/eatmoreapple/openwechat"
	"log"
	"os"
	loader "sf-bot/handler/load"
	"strings"
)

type user struct {
	id        string
	friend    *openwechat.Friend
	gptEnable bool
}

type userHandler struct {
	users   map[string]*user
	bdInfos map[string]BDInfo
}

type BDInfo struct {
	Id        string `json:"id"`
	NickName  string `json:"nickName"`
	GptEnable bool   `json:"gptEnable"`
}

const bdFile = "./user_handler_bds.json"

var users *userHandler

func (h *userHandler) BD() {
	if len(h.bdInfos) > 0 {
		if marshal, err := json.Marshal(h.bdInfos); err == nil {
			err := os.WriteFile(bdFile, marshal, 0775)
			if err != nil {
				log.Println("brush disk error ", err)
			}
		}
	}
}

func (h *userHandler) loadBD() {
	_, err := os.Stat(bdFile)
	if err == nil {
		if file, err := os.ReadFile(bdFile); err == nil {
			i := new(map[string]BDInfo)
			json.Unmarshal(file, i)
			infos := *i
			h.bdInfos = infos
			log.Println("load bds ", infos)
		}
	}
}

func (h *userHandler) Recv(ctx *openwechat.MessageContext) {
	if ctx.IsText() {

		id := ""
		if sender, err := ctx.Sender(); err == nil {
			friend, _ := sender.AsFriend()
			id = friend.PYQuanPin + friend.RemarkPYQuanPin
			h.assemblyUser(id, friend)

		}

		if flag := h.command(id, ctx.Content, func(msg string) { ctx.ReplyText(msg) }); flag {
			return
		}

	}
}

func (h *userHandler) command(id, content string, f func(msg string)) bool {
	if strings.Contains(content, "#gpt-enable#") {
		h.users[id].gptEnable = true
		info := h.bdInfos[id]
		info.GptEnable = true
		h.bdInfos[id] = info
		f(h.ToString())
		return true
	}
	if strings.Contains(content, "#gpt-disable#") {
		h.users[id].gptEnable = false
		info := h.bdInfos[id]
		info.GptEnable = false
		h.bdInfos[id] = info
		f(h.ToString())
		return true
	}
	//todo 2023/6/8 lamkeizyi -

	return false

}

func (h *userHandler) assemblyUser(id string, friend *openwechat.Friend) {
	u := h.users[id]
	if u != nil {
		return
	}
	u2 := &user{id, friend, false}
	h.users[id] = u2

	//如果从本地临时文件获取到配置，则采用文件中的配置
	info := h.bdInfos[id]
	if len(info.Id) > 0 {
		u2.gptEnable = info.GptEnable
	} else {
		bdInfo := BDInfo{id, friend.NickName, false}
		h.bdInfos[id] = bdInfo
	}
}

func (h *userHandler) ToString() string {
	//todo 2023/6/8 lamkeizyi -

	return ""
}

func CreateUserHandler() *userHandler {
	if loader.LoadBool("user.enable") {
		u := &userHandler{}
		u.loadBD()
		users = u
		return u

	}
	return nil
}
