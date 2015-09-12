package http

import (
	// Standard library
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// Internal packages
	"github.com/deuill/sigil/lib/config"
	"github.com/deuill/sigil/lib/engine"
	"github.com/deuill/sigil/srv"
)

type Config struct {
	Base  string
	Hosts string
	Index string
	Root  string
}

type Service struct {
	port  *string
	base  *string
	hosts map[string]*Config
}

func (s *Service) Init() error {
	// Read configuration for hosts, if any exist.
	if err := s.setup(); err != nil {
		return err
	}

	// Listen and serve on predefined HTTP port.
	ln, err := net.Listen("tcp", net.JoinHostPort("", *s.port))
	if err != nil {
		return err
	}

	http.HandleFunc("/", s.handle)
	go http.Serve(ln, nil)

	return nil
}

func (s *Service) handle(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)

	// Return 404 if configuration does not exist for the host.
	if _, exists := s.hosts[host]; !exists {
		s.fail(w, 404)
		return
	}

	// Handle panics gracefully (by showing a 500 error).
	defer func() {
		if r := recover(); r != nil {
			s.fail(w, 500)
		}
	}()

	// Attempt to load file pointed to by URL path, moving on to the default engine handler for the
	// specified host if no file exists (or is otherwise inaccessible) for that location.
	name := s.hosts[host].Base + "/" + s.hosts[host].Root + "/" + r.URL.Path
	if fi, err := os.Stat(name); err != nil || fi.IsDir() {
		index := s.hosts[host].Base + "/" + s.hosts[host].Index
		if err := engine.Handle(w, index); err != nil {
			s.fail(w, 500)
			return
		}
	} else {
		http.ServeFile(w, r, name)
	}
}

func (s *Service) fail(w http.ResponseWriter, code int) {
	status := http.StatusText(code)

	// Return error 500 if error code passed is undefined.
	if status == "" {
		s.fail(w, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(`{"status": "` + status + `"}`))
}

func (s *Service) setup() error {
	sites, err := filepath.Glob(*s.base + "/*")
	if err != nil {
		return err
	}

	for _, base := range sites {
		// Skip unreadable or non-directory matches.
		fi, err := os.Stat(base)
		if err != nil || !fi.IsDir() {
			continue
		}

		// Skip directories with missing configuration.
		if _, err = os.Stat(base + "/config/sigil.conf"); err != nil {
			continue
		}

		// Load configuration file, and attach configuration variables to hosts configuration.
		conf, err := config.Load(base + "/config/sigil.conf")
		if err != nil {
			return err
		}

		c := &Config{Base: base}
		if err = conf.Unpack(c); err != nil {
			return err
		}

		// Attach host definition for all hostnames defined in current site.
		for _, h := range strings.Fields(c.Hosts) {
			// Do not overwrite already existing host.
			if _, exists := s.hosts[h]; exists {
				return fmt.Errorf("Configuration for hostname '%s' already exists", h)
			}

			s.hosts[h] = c
		}
	}

	return nil
}

func init() {
	flags := flag.NewFlagSet("http", flag.ContinueOnError)
	service := &Service{
		port:  flags.String("port", "80", ""),
		base:  flags.String("base", "/srv/http", ""),
		hosts: make(map[string]*Config),
	}

	srv.Register("http", service, flags)
}
