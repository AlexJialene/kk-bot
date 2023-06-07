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
	bdInfos map[string]*BDInfo
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
			i := new(map[string]*BDInfo)
			json.Unmarshal(file, i)
			h.bdInfos = *i
			log.Println("load bds ", h.bdInfos)
		}
	}
}

func (h *userHandler) Recv(ctx *openwechat.MessageContext) {
	if ctx.IsText() {
		log.Printf("recv friend message = %s", ctx.Content)

		id := ""
		if sender, err := ctx.Sender(); err == nil {
			friend, _ := sender.AsFriend()
			id = friend.PYQuanPin + friend.RemarkPYQuanPin
			h.assemblyUser(id, friend)
		}

		if flag := h.command(id, ctx.Content, func(msg string) { ctx.ReplyText(msg) }); flag {
			return
		}

		if h.bdInfos[id].GptEnable {
			ask := agent.AskAgent().Ask(id, ctx.Content)
			if _, err := ctx.ReplyText(ask); err != nil {
				log.Println("reply text error ", err)
			}
		}
	}
}

func (h *userHandler) command(id, content string, f func(msg string)) bool {
	if strings.Contains(content, "#gpt-enable#") {
		h.users[id].gptEnable = true
		h.bdInfos[id].GptEnable = true
		f(h.ToString(id))
		return true
	}
	if strings.Contains(content, "#gpt-disable#") {
		h.users[id].gptEnable = false
		h.bdInfos[id].GptEnable = false
		f(h.ToString(id))
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
	if info != nil {
		u2.gptEnable = info.GptEnable
	} else {
		bdInfo := &BDInfo{id, friend.NickName, false}
		h.bdInfos[id] = bdInfo
	}
}

func (h *userHandler) ToString(id string) string {
	var result = ""
	result += "======setting======\n"
	if h.bdInfos[id].GptEnable {
		result += "gpt-enable = true \n"
	}
	if !h.bdInfos[id].GptEnable {
		result += "gpt-enable = false \n"
	}
	result += "======setting======"
	return result
}

func CreateUserHandler() *userHandler {
	if loader.LoadBool("user.enable") {
		u := &userHandler{
			make(map[string]*user),
			make(map[string]*BDInfo),
		}
		u.loadBD()
		users = u
		return u

	}
	return nil
}
