package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/rakyll/globalconf"
	"github.com/thoughtmonster/sigil/serv"
)

func main() {
	opts := &globalconf.Options{EnvPrefix: "SIGIL_", Filename: "conf/sigil.conf"}
	conf, err := globalconf.NewWithOptions(opts)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	conf.ParseAll()

	fmt.Print("Starting server... ")

	err = serv.Init()
	if err != nil {
		fmt.Printf("error initializing services:\n%s\n", err)
		os.Exit(1)
	}

	fmt.Println("started successfully.")

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
