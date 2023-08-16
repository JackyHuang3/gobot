package webhook

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type User struct {
	ID string `json:"id"`
}

type Message struct {
	ID     string `json:"mid"`
	Text   string `json:"text"`
	AppID  uint64 `json:"app_id"`
	IsEcho bool   `json:"is_echo"`
}

func MsgHandle(data []byte) (bool, error) {
	if !gjson.Get(string(data), "message").Exists() {
		return false, nil
	}
	messaging := struct {
		Sender    User    `json:"sender"`
		Recipient User    `json:"recipient"`
		Message   Message `json:"message"`
		Timestamp uint64  `json:"timestamp"`
	}{}
	LogPrintf("MsgHandle receive data: %s", data)
	if err := json.Unmarshal([]byte(data), &messaging); err != nil {
		return true, err
	}
	LogPrintf("MsgHandle parsed Messaging: %+v", messaging)
	return true, messaging.Message.Receive(messaging.Sender.ID)
}

func (p *Message) Receive(senderID string) error {
	go p.Handle(senderID)
	return nil
}

func (p *Message) Handle(senderID string) {
	// TODO: 处理用户消息
	if p.Text == "feedback" {
		resp := &Response{
			Recipient: senderID,
		}
		resp.Send(resp.FeedbackTemplate("gobot", "gobot"))
	} else {
		resp := &Response{
			Recipient: senderID,
		}
		resp.Send(resp.TextTemplate(p.Text))
	}
}
