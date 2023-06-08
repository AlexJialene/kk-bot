package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron"
	"log"
	"os"
	"reflect"
	loader "sf-bot/handler/load"
	"strconv"
	"strings"
	"time"
)

type user struct {
	id     string
	friend *openwechat.Friend

	//废弃-未用到
	gptEnable bool
}

type userHandler struct {
	users   map[string]*user
	bdInfos map[string]*BDInfo
}

type BDInfo struct {
	Id            string `json:"id"`
	NickName      string `json:"nickName"`
	GptEnable     bool   `json:"gptEnable"`
	ZhihuPush     bool   `json:"zhihuPush"`
	WeiboPush     bool   `json:"weiboPush"`
	NewsPush      bool   `json:"newsPush"`
	ZhihuPushHour int    `json:"zhihuPushHour"`
	WeiboPushHour int    `json:"weiboPushHour"`
	NewsPushHour  int    `json:"newsPushHour"`
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
	pindex := strings.Index(content, "#")
	lindex := strings.LastIndex(content, "#")
	if pindex != -1 && lindex != -1 && pindex != lindex {
		s := content[pindex+1 : lindex]
		if strings.IndexAny(s, "-") == -1 {
			f(h.ToString(id))
			return true
		}
		split := strings.Split(s, "-")
		i := len(split)
		if i == 2 || i == 3 {
			info := h.bdInfos[id]
			of := reflect.ValueOf(info)
			elem := of.Elem()
			name := elem.FieldByName(split[0])
			if split[1] == "enable" {
				name.SetBool(true)
				if i != 3 {
					f(h.ToString(id))
					return true
				}
			}
			if split[1] == "disable" {
				name.SetBool(false)
				if i != 3 {
					f(h.ToString(id))
					return true
				}
			}
			if i == 3 {
				if atoi, err := strconv.Atoi(split[0]); err == nil {
					hourName := elem.FieldByName(split[0] + "Hour")
					hourName.SetInt(int64(atoi))
					f(h.ToString(id))
					return true
				}
			}

		}
	}

	//if strings.Contains(content, "#gpt-enable#") {
	//	h.users[id].gptEnable = true
	//	h.bdInfos[id].GptEnable = true
	//	f(h.ToString(id))
	//	return true
	//}
	//if strings.Contains(content, "#gpt-disable#") {
	//	h.users[id].gptEnable = false
	//	h.bdInfos[id].GptEnable = false
	//	f(h.ToString(id))
	//	return true
	//}

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
		bdInfo := &BDInfo{Id: id, NickName: friend.NickName, GptEnable: false}
		h.bdInfos[id] = bdInfo
	}
}

func (h *userHandler) ToString(id string) string {
	var result = ""
	result += "------setting------\n"

	info := h.bdInfos[id]
	of := reflect.ValueOf(info)
	typeOf := reflect.TypeOf(info)

	typeField := of.Elem().NumField()
	for i := 0; i < typeField; i++ {
		name := typeOf.Elem().Field(i).Name
		field := of.Elem().Field(i)
		v := ""
		if field.Kind().String() == "string" {
			v = field.String()
		}
		if field.Kind().String() == "int" {
			i2 := field.Int()
			v = strconv.FormatInt(i2, 10)
		}
		if field.Kind().String() == "bool" {
			v = fmt.Sprintf("%t", field.Bool())
		}
		if name != "Id" {
			result += name + " : " + v + "\n"
		}
	}
	result += "------setting------"
	return result
}

func (h *userHandler) timePush() {
	go h.BD()
	hour := time.Now().Hour()
	if hour == 0 {
		return
	}
	if len(h.bdInfos) > 0 {
		for _, v := range h.bdInfos {

			if v.ZhihuPush && v.ZhihuPushHour == hour {
				//h.users[v.Id].friend.SendText();
			}
			if v.WeiboPush && v.WeiboPushHour == hour {

			}
			if v.NewsPush && v.NewsPushHour == hour {

			}
		}

	}
}

func CreateUserHandler() *userHandler {
	if loader.LoadBool("user.enable") {
		u := &userHandler{
			make(map[string]*user),
			make(map[string]*BDInfo),
		}
		u.loadBD()
		users = u

		go func() {
			c := cron.New()
			c.AddFunc("0 1 0/1 * * ?", u.timePush)
			c.Start()
			select {}
		}()
		return u

	}
	return nil
}
