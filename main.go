package main

import (
	// Standard library
	"fmt"
	"os"
	"os/signal"
	"strings"

	// Internal packages
	"github.com/deuill/sigil/srv"

	// External packages
	"github.com/rakyll/globalconf"
)

func setup(name string) error {
	var (
		conf *globalconf.GlobalConf
		file string
		err  error
	)

	// Check for environment variable override, and return error if override is set but file is
	// inaccessible. Otherwise, try for the default global location (in the "/etc" directory).
	if file := os.Getenv(strings.ToUpper(name) + "_CONFIG"); file != "" {
		if _, err = os.Stat(file); err != nil {
			return err
		}
	} else {
		file = "/etc/" + name + "/" + name + ".conf"
		if _, err = os.Stat(file); err != nil {
			file = ""
		}
	}

	// Load from specific configuration file if set, or use local configuration file as a fallback.
	if file != "" {
		options := &globalconf.Options{Filename: file, EnvPrefix: ""}
		if conf, err = globalconf.NewWithOptions(options); err != nil {
			return err
		}
	} else if conf, err = globalconf.New(name); err != nil {
		return err
	}

	conf.EnvPrefix = strings.ToUpper(name) + "_"
	conf.ParseAll()

	return nil
}

func main() {
	var err error

	if err = setup("sigil"); err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	fmt.Print("Starting server... ")

	// Initialize HTTP and attached services.
	err = srv.Init()
	if err != nil {
		fmt.Printf("error initializing services:\n%s\n", err)
		os.Exit(1)
	}

	fmt.Println("done.")

	// Listen for and terminate Sigil on SIGKILL or SIGINT signals.
	sigStop := make(chan os.Signal)
	signal.Notify(sigStop, os.Interrupt, os.Kill)

	select {
	case <-sigStop:
		fmt.Println("\rShutting down server...")
	}

	os.Exit(0)
}
