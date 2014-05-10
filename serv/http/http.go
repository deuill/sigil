package http

import (
	"flag"

	"github.com/thoughtmonster/crowley/serv"
)

type HTTPService struct {
	host *string
	port *string
}

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
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	h := &HTTPService{
		host: fs.String("host", "127.0.0.1", ""),
		port: fs.String("port", "8080", ""),
	}

	serv.Register("http", h, fs)
}
