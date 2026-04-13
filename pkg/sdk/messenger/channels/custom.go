package channels

import "message-nest/pkg/message"

type CustomChannel struct{ *BaseChannel }

func NewCustomChannel() Channel {
	return &CustomChannel{NewBaseChannel(ChannelCustom, []string{FormatTypeText})}
}

func (c *CustomChannel) Send(config ChannelConfig, msg *Message) (*Result, error) {
	webhook := config.GetString("webhook")
	header := config.GetString("header")
	body := config.GetString("body")

	if webhook == "" {
		return SendError("custom config missing: webhook is required"), nil
	}

	_, formattedContent := c.FormatContent(msg)
	cli := message.CustomWebhook{}

	// 替换 body 模板中的 TEXT 占位符
	bodyStr := body
	if bodyStr != "" {
		bodyStr = replaceBodyPlaceholder(bodyStr, formattedContent)
	} else {
		bodyStr = formattedContent
	}

	headers := cli.ParseHeaders(header)

	res, err := cli.Request(webhook, bodyStr, headers)
	if err != nil {
		return ErrorResult(string(res), err), nil
	}
	return SuccessResult(string(res)), nil
}
