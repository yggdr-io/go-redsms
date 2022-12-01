package redsms

import "context"

type MessageService service

func (s *MessageService) Send(ctx context.Context, msg *Message) (*SendMessageReport, *Response, error) {
	u := "message"
	req, err := s.client.NewRequest("POST", u, msg)
	if err != nil {
		return nil, nil, err
	}

	var report SendMessageReport
	resp, err := s.client.Do(ctx, req, &report)
	if err != nil {
		return nil, resp, err
	}

	return &report, resp, nil
}

type MessageRoute string

const (
	MessageRouteSMS MessageRoute = "sms"
)

type Message struct {
	From  string       `json:"from"`
	To    string       `json:"to"`
	Text  string       `json:"text"`
	Route MessageRoute `json:"route"`
}

type SendMessageReport struct {
	Items   []*SendMessagReportItem   `json:"items"`
	Errors  []*SendMessageReportError `json:"errors"`
	Success bool                      `json:"success"`
}

type SendMessagReportItem struct {
	UUID string `json:"uuid"`
	To   string `json:"to"`
}

type SendMessageReportError struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
