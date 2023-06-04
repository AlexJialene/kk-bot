package gpt

import (
	"github.com/google/uuid"
	"log"
	"time"
)

// connect to pandora
// pandora: https://github.com/pengzhile/pandora
var token string

var senders map[string]*Sender

type Sender struct {
	id              string
	messageId       string
	parentMessageId string
	conversationId  string
	createTime      int64
}

func GetSender(id string) *Sender {
	sender := senders[id]
	if sender != nil {
		sender.messageId = uuid.NewString()
		sender.createTime = time.Now().Unix()
	}
	return sender
}

func newSender(id string) *Sender {
	if len(senders) == 0 {
		senders = make(map[string]*Sender)
	}
	sender := &Sender{
		id,
		uuid.New().String(),
		uuid.New().String(),
		"",
		time.Now().Unix(),
	}
	senders[id] = sender
	return sender
}

func setMessageId(id, conversationId, parentMessageId string) {
	sender := senders[id]
	if sender != nil {
		sender.conversationId = conversationId
		sender.parentMessageId = parentMessageId
	}
}

func init() {
	go func() {
		if len(senders) == 0 {
			senders = make(map[string]*Sender)
		}

		//创建一个维护上下文的监听线程
		//有效期60s
		for {
			curr := time.Now().Unix()
			for _, sender := range senders {
				if curr-sender.createTime > 60 {
					log.Printf("%s was expired ", sender.id)
					delete(senders, sender.id)
				}
			}
			time.Sleep(5000)
		}
	}()
}
