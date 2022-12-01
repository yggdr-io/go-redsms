package redsms

import "context"

type ClientService service

func (s *ClientService) GetInfo(ctx context.Context) (*ClientInfo, *Response, error) {
	u := "client/info"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var respAPI struct {
		Info *ClientInfo `json:"info"`
	}
	resp, err := s.client.Do(ctx, req, &respAPI)
	if err != nil {
		return nil, nil, err
	}

	return respAPI.Info, resp, nil
}

type ClientInfo struct {
	Login     string   `json:"login"`
	Balance   float64  `json:"balance"`
	Active    bool     `json:"active"`
	Overdraft *float64 `json:"overdraft,omitempty"`
}
