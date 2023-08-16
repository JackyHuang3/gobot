package webhook

import (
	"fmt"
	"gobot/global"
	"strings"
)

type Response struct {
	Recipient string
}

func (p *Response) Send(sendJson string) error {
	sendAPIUrl := fmt.Sprintf("%s/%s/messages?access_token=%s", vConfig.Get(global.C_ApiURL), vConfig.Get(global.C_PageID), vConfig.Get(global.C_AccessToken))
	LogPrintf("before send message, sendAPIUrl: %s, sendJson: %s", sendAPIUrl, sendJson)

	return WebPostByFunc(sendAPIUrl, strings.NewReader(sendJson), func(buf []byte) error {
		LogPrintf("messaging send response: %s", string(buf))
		return nil
	})
}

func (p *Response) TextTemplate(text string) string {
	output := `{"recipient":{"id":"recipient_id"},"messaging_type":"RESPONSE","message":{"text":"text_content"}}`
	output = strings.Replace(output, "recipient_id", p.Recipient, -1)
	output = strings.Replace(output, "text_content", text, -1)
	return output
}

func (p *Response) FeedbackTemplate(productID, productName string) string {
	output := `{"recipient":{"id":"recipient_id"},"message":{"attachment":{"type":"template","payload":{"template_type":"customer_feedback","title":"Rate your experience with Original Coast Clothing.","subtitle":"Let product_name know how they are doing by answering two questions","button_title":"Rate Experience","feedback_screens":[{"questions":[{"id":"product_id","type":"csat","title":"How would you rate your experience with product_name?","score_label":"neg_pos","score_option":"five_stars","follow_up":{"type":"free_form","placeholder":"Give additional feedback"}}]}],"business_privacy":{"url":"https://www.facebook.com"},"expires_in_days":3}}}}`
	output = strings.Replace(output, "recipient_id", p.Recipient, -1)
	output = strings.Replace(output, "product_name", productName, -1)
	output = strings.Replace(output, "product_id", productID, -1)
	return output
}
