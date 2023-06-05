package handler

import (
	"fmt"
	"github.com/sf-bot/gpt"
	"log"
	loader "sf-bot/handler/load"
)

type AgentFactory struct {
	gpt   DataChannel
	redis DataChannel
	cuz   bool
	c     uint8
}

func InitAgent() *AgentFactory {
	var g, r DataChannel
	gptEnable := loader.LoadBool("gpt.enable")
	if gptEnable {
		host := loader.Load("gpt.host")
		g = Gpt{host: host, gptClient: gpt.CreateGptHost(host)}
		log.Printf("gpt host = %s \n ", host)
	}

	redisEnable := loader.LoadBool("gpt.enable")
	if redisEnable {
		redis := loader.Load("redis.host")
		r = Redis{host: redis}
		log.Printf("redis host = %s \n", redis)
	}

	return &AgentFactory{
		g,
		r,
		loader.LoadBool("commons.cuz"),
		0,
	}
}

func (a *AgentFactory) GetAgent(string) DataChannel {

	return nil
}

func (a *AgentFactory) AskAgent() DataChannel {
	if a.cuz {
		//todo 2023/6/4 lamkeizyi -
	}
	if a.gpt != nil {
		return a.gpt
	}
	if a.redis != nil {
		return a.redis
	}
	return DefaultChannel{}
}

type DataChannel interface {
	Ask(id, question string) string
}

type Gpt struct {
	//http://127.0.0.1:19090
	host      string
	gptClient *gpt.GptClient
}

type Redis struct {
	host string
	port string
	auth string
	db   int32
}

func (g Gpt) Ask(id, question string) string {
	fmt.Println(id)
	fmt.Println(question)
	talk := g.gptClient.Talk(id, question)
	return talk

}

func (g Redis) Ask(id, question string) string {
	//todo 2023/6/4 lamkeizyi -
	return ""
}

type DefaultChannel struct {
}

func (g DefaultChannel) Ask(id, question string) string {
	return "system not setting any questions"
}
