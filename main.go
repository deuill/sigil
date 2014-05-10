package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/rakyll/globalconf"
	"github.com/thoughtmonster/crowley/serv"

	_ "github.com/thoughtmonster/crowley/serv/http"
)

func main() {
	opts := &globalconf.Options{EnvPrefix: "CROWLEY_"}
	conf, err := globalconf.NewWithOptions(opts)
	if err != nil {
		log.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	conf.ParseAll()

	err = serv.Init()
	if err != nil {
		log.Println("Error initializing services:", err)
		os.Exit(1)
	}

	log.Println("Server start")

	sigStop := make(chan os.Signal)
	signal.Notify(sigStop, os.Interrupt, os.Kill)

	select {
	case <-sigStop:
		log.Println("Shutting down server...")

		errs := serv.Shutdown()
		if errs != nil {
			log.Println("The following services failed to shut down cleanly:")
			for _, err = range errs {
				log.Println(err)
			}

			log.Println("The environment might be in an unclean state")
			os.Exit(2)
		}
	}

	os.Exit(0)
}
