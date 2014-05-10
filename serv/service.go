package serv

import (
	"flag"
	"fmt"

	"github.com/rakyll/globalconf"
)

type Service interface {
	Setup() error
	Start() error
	Stop() error
}

var services map[string]Service

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

func Init() error {
	var err error

	for _, s := range services {
		if err = s.Setup(); err != nil {
			return err
		}

		if err = s.Start(); err != nil {
			return err
		}
	}

	return nil
}

func Shutdown() []error {
	var err error
	var errs []error

	for name, s := range services {
		if err = s.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", name, err))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func init() {
	services = make(map[string]Service)
}
