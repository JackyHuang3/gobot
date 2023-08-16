package webhook

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

func FeedbackHandle(data []byte) (bool, error) {
	if !gjson.Get(string(data), "messaging_feedback").Exists() {
		return false, nil
	}
	messaging := struct {
		Sender    User     `json:"sender"`
		Receiver  User     `json:"recipient"`
		Feedback  Feedback `json:"messaging_feedback"`
		Timestamp uint64   `json:"timestamp"`
	}{}
	LogPrintf("FeedbackHandle receive data: %s", data)
	if err := json.Unmarshal([]byte(data), &messaging); err != nil {
		return true, err
	}
	LogPrintf("FeedbackHandle parsed Messaging: %+v", messaging)
	return true, messaging.Feedback.Receive(messaging.Sender.ID)
}

type FeedbackFollowUp struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type FeedbackMsg struct {
	Type     string           `json:"type"`
	Payload  string           `json:"payload"`
	FollowUP FeedbackFollowUp `json:"follow_up"`
}

type Screen struct {
	ID        uint64                       `json:"screen_id"`
	Questions map[string]*FeedbackFollowUp `json:"questions"`
}

type Feedback struct {
	Screens []Screen `json:"feedback_screens"`
}

func (p *Feedback) Receive(senderID string) error {
	go p.Handle(senderID)
	return nil
}

func (p *Feedback) Handle(senderID string) {
	// store feedback
	fdBuf, err := json.Marshal(p.Screens)
	if err != nil {
		LogPrintf("Feedback Handle failed, detail: %s", err.Error())
		return
	}
	vStore.Set(senderID, fdBuf)

	// analyze
	responseText := p.analyze(senderID)

	// send notify
	resp := &Response{
		Recipient: senderID,
	}
	resp.Send(resp.TextTemplate(responseText))
}

func (p *Feedback) analyze(senderID string) string {
	for _, screen := range p.Screens {
		for _, question := range screen.Questions {
			if question.Type != "CSAT" {
				return "Thanks for your feedback"
			}
			switch question.Payload {
			case "5":
				return fmt.Sprintf("Thank you for your %s star feedback, we will continue to improve and optimize", question.Payload)
			case "4":
				return fmt.Sprintf("Thank you for your %s star feedback, we will continue to improve and optimize", question.Payload)
			case "3":
				return fmt.Sprintf("Thank you for your %s star feedback, we will continue to improve and optimize", question.Payload)
			case "2":
				return fmt.Sprintf("Thank you for your %s star feedback, we will continue to improve and optimize", question.Payload)
			case "1":
				return fmt.Sprintf("Thank you for your %s star feedback, we will continue to improve and optimize", question.Payload)
			}
		}
	}
	return ""
}
