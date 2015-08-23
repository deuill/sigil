package srv

import (
	// Standard library
	"flag"
	"fmt"

	// External packages
	"github.com/rakyll/globalconf"
)

type Service interface {
	Init() error
}

// A map of all registered services.
var services map[string]Service

// Register attaches services, to be initialized in a future date.
func Register(name string, rcvr Service, fs *flag.FlagSet) error {
	if _, exists := services[name]; exists {
		return fmt.Errorf("Service '%s' already exists, refusing to overwrite", name)
	}

	if fs != nil {
		globalconf.Register(name, fs)
	}

	services[name] = rcvr
	return nil
}

// Initialize all registered services, using the Service.Init method.
func Init() error {
	var err error

	for name, s := range services {
		if err = s.Init(); err != nil {
			return fmt.Errorf("[%s]: %s", name, err)
		}
	}

	return nil
}

func init() {
	services = make(map[string]Service)
}
