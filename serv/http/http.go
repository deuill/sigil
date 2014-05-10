package http

import (
	"github.com/thoughtmonster/crowley/serv"
)

type HTTPService struct{}

func (s *HTTPService) Setup() error {
	return nil
}

func (s *HTTPService) Start() error {
	return nil
}

func (s *HTTPService) Stop() error {
	return nil
}

func init() {
	serv.Register("http", &HTTPService{})
}
