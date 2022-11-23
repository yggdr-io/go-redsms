package redsms

import "context"

type MessageService service

func (s *MessageService) Send(ctx context.Context, msg *Message) (*SendMessageResponse, *Response, error) {
	u := "message"
	req, err := s.client.NewRequest("POST", u, msg)
	if err != nil {
		return nil, nil, err
	}

	var respAPI SendMessageResponse
	resp, err := s.client.Do(ctx, req, &respAPI)
	if err != nil {
		return nil, resp, err
	}

	return &respAPI, resp, nil
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

type SendMessageResponse struct {
	Items   []*SendMessageItem  `json:"items"`
	Errors  []*SendMessageError `json:"errors"`
	Success bool                `json:"success"`
}

type SendMessageItem struct {
	UUID string `json:"uuid"`
	To   string `json:"to"`
}

type SendMessageError struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
