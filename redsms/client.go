package redsms

import "context"

type ClientService service

func (s *ClientService) GetInfo(ctx context.Context) (*Response, error) {
	u := "client/info"
	_, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
