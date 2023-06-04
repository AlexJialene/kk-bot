package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	talk       = "/api/conversation/talk"
	gen_title  = "/api/conversation/gen_title/"
	regenerate = "/api/conversation/regenerate"

	//keep continue saying
	goon = "/api/conversation/goon"
)

type GptClient struct {
	host string
}

func CreateGptHost(h string) *GptClient {
	return &GptClient{h}
}

func matchSender(id string) *Sender {
	sender := GetSender(id)
	if sender != nil {
		return sender
	} else {
		return newSender(id)
	}
}

func (gptClient *GptClient) Talk(id, question string) string {
	sender := matchSender(id)
	url := gptClient.host + talk
	talkVO := TalkVO{
		question,
		"gpt-3.5-turbo",
		sender.messageId,
		sender.parentMessageId,
		sender.conversationId,
		false,
	}
	log.Printf("req_message_id = %s \n parnet_message_id= %s \n conversation_id = %s",
		sender.messageId,
		sender.parentMessageId,
		sender.conversationId,
	)

	marshal, _ := json.Marshal(talkVO)
	reader := strings.NewReader(string(marshal))
	post, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Println(err)
		return "the gpt talking has error "
	}
	defer post.Body.Close()
	all, err := io.ReadAll(post.Body)
	if err != nil {
		return "read body error "
	}

	t := &TalkRespVO{}
	_ = json.Unmarshal(all, t)
	conversation_id := t.Conversation_id
	message_id := t.Message.Id

	log.Printf("response:\n message_id = %s \n conversation_id = %s", message_id, conversation_id)

	if len(sender.conversationId) == 0 {
		log.Println("request generate title ...")
		defer gptClient.genTitle(conversation_id, "gpt-3.5-turbo", message_id)
	}
	setMessageId(id, conversation_id, message_id)
	parts := t.Message.Content.Parts
	answer := ""
	for _, content := range parts {
		answer = answer + content
	}
	return answer
}

func (gptClient *GptClient) ReTalk() {

}

func (gptClient *GptClient) genTitle(conversationId, model, messageId string) string {
	url := gptClient.host + gen_title + conversationId
	j := Json{model, messageId}
	marshal, _ := json.Marshal(j)
	fmt.Println(string(marshal))
	reader := strings.NewReader(string(marshal))
	post, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Println(err)
		return "the gpt talking has error "
	}

	defer post.Body.Close()
	body, _ := io.ReadAll(post.Body)
	fmt.Println(string(body))
	t := &Title{}
	_ = json.Unmarshal(body, t)
	log.Printf("generate title  = %s", t.Title)
	return t.Title

}
