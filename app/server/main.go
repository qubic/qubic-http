package main

import (
	"fmt"
	"github.com/0xluk/go-qubic"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"qubic-api-sidecar/app/server/handlers"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
)

const prefix = "QUBIC_API_SIDECAR"

func main() {
	if err := run(); err != nil {
		log.Fatalf("main: exited with error: %s", err.Error())
	}
}

func run() error {
	var cfg struct {
		Web struct {
			Host            string        `conf:"default:0.0.0.0:8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Qubic struct {
			NodeIP   string `conf:"default:65.21.10.217"`
			NodePort string `conf:"default:21841"`
		}
	}

	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	qubicClient, err := qubic.NewClient(cfg.Qubic.NodeIP, cfg.Qubic.NodePort)
	if err != nil {
		return errors.Wrap(err, "creating qubic client")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", cfg.Web.Host)
		serverErrors <- http.ListenAndServe(cfg.Web.Host, handlers.New(qubicClient))
	}()

	select {
	case err := <-serverErrors:
		return err

	case <-shutdown:
		return errors.New("shutdown initialized")
	}
}
