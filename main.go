package main

import (
	// Standard library
	"fmt"
	"os"
	"os/signal"

	// Internal packages
	"github.com/deuill/sigil/serv"

	// External packages
	"github.com/rakyll/globalconf"
)

// Entry point for Sigil, this sets up global configuration and starts internal services.
func main() {
	conf, err := globalconf.New("sigil")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	// Initialize configuration, reading from environment variables using a 'SIGIL_' prefix first,
	// then moving to a static configuration file, usually located in ~/.config/sigil/config.ini.
	conf.EnvPrefix = "SIGIL_"
	conf.ParseAll()

	fmt.Print("Starting server... ")

	// Initialize HTTP and attached services.
	err = serv.Init()
	if err != nil {
		fmt.Printf("error initializing services:\n%s\n", err)
		os.Exit(1)
	}

	fmt.Println("started successfully.")

	// Listen for and terminate Sigil on SIGKILL or SIGINT signals.
	sigStop := make(chan os.Signal)
	signal.Notify(sigStop, os.Interrupt, os.Kill)

	select {
	case <-sigStop:
		fmt.Println("Shutting down server...")

		errs := serv.Shutdown()
		if errs != nil {
			fmt.Println("The following services failed to shut down cleanly:")
			for _, err = range errs {
				fmt.Println(err)
			}

			fmt.Println("The environment might be in an unclean state")
			os.Exit(2)
		}
	}

	os.Exit(0)
}
