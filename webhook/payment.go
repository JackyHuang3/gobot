package webhook

func PaymentHandle(id, actions string) (bool, error) {
	// Webhooks 更新只会通知您特定支付（由 id 字段标识）已发生更改。收到更新后，您接下来需要查询图谱 API 以了解交易详情，从而恰当地处理更改
	return false, nil
}

type Payment struct {
}

func (p *Payment) Receive(senderID string) error {
	go p.Handle(senderID)
	return nil
}

func (p *Payment) Handle(senderID string) {
	// send notify
}
