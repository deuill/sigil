package http

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ini "github.com/rakyll/goini"
	"github.com/thoughtmonster/crowley/serv"
)

type HTTPService struct {
	port  *string
	base  *string
	conf  *string
	sites map[string]ini.Dict
}

func (s *HTTPService) handleRequest(w http.ResponseWriter, r *http.Request) {
	base := *s.base + "/" + r.Host

	// Override global variables with site configuration, if any exists.
	if t, exists := s.sites[r.Host]; exists {
		base, _ = t.GetString("", "base")
	}

	// Return 404 error if the main controller isn't found.
	if _, err := os.Stat(base + "/main.js"); err != nil {
		s.handleError(404, w, r)
		return
	}

	fmt.Fprintf(w, "Loaded %s!", base+"/main.js")
}

func (s *HTTPService) handleError(code int, w http.ResponseWriter, r *http.Request) {
	// Return error 500 if error code passed is undefined.
	if _, exists := errorCodes[code]; !exists {
		s.handleError(500, w, r)
		return
	}

	w.WriteHeader(code)
	tpl, _ := template.New("error").Parse(errorTemplate)
	tpl.Execute(w, errorCodes[code])
}

func (s *HTTPService) loadConf() error {
	files, err := filepath.Glob(*s.conf + "/*.conf")
	if err != nil {
		return err
	}

	for _, name := range files {
		t, err := ini.Load(name)
		if err != nil {
			return err
		}

		h, ok := t.GetString("", "hostname")
		if !ok {
			return fmt.Errorf("Could not find 'hostname' value in configuration for '%s'", name)
		}

		if _, ok = s.sites[h]; ok {
			return fmt.Errorf("Configuration for hostname '%s' already exists", h)
		}

		s.sites[h] = t
	}

	return nil
}

func (s *HTTPService) Start() error {
	// Read configuration for sites, if any exists.
	if err := s.loadConf(); err != nil {
		return err
	}

	// Start HTTP server, sending any errors back to the 'result' channel.
	result := make(chan error)
	go func() {
		http.HandleFunc("/", s.handleRequest)
		result <- http.ListenAndServe(":"+*s.port, nil)
	}()

	// Allow for a 500 millisecond timeout before this function returns, in order to catch any
	// errors that might be emitted by the HTTP server goroutine.
	timeout := make(chan bool)
	go func() {
		time.Sleep(500 * time.Millisecond)
		timeout <- true
	}()

	select {
	case err := <-result:
		// The HTTP server has returned before the timeout, which means that an error has occured.
		return err
	case <-timeout:
		// A timeout has occurred and no error has been received by the server, which probably
		// means that everything went OK.
		return nil
	}
}

func (s *HTTPService) Stop() error {
	return nil
}

func init() {
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	h := &HTTPService{
		port:  fs.String("port", "80", ""),
		base:  fs.String("base", "/srv/http", ""),
		conf:  fs.String("conf", "/etc/crowley/sites.d", ""),
		sites: make(map[string]ini.Dict),
	}

	serv.Register("http", h, fs)
}
