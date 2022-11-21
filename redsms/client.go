package redsms

import "context"

type ClientService service

func (s *ClientService) GetInfo(ctx context.Context) (*Response, error) {
	u := "client/info"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
