package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qubic/qubic-http/app/server/handlers"
	"github.com/qubic/qubic-http/external/opensearch"
	"github.com/qubic/qubic-http/foundation/nodes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
)

const prefix = "QUBIC_API_SIDECAR"

func main() {
	log := log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		log.Fatalf("main: exited with error: %s", err.Error())
	}
}

func run(log *log.Logger) error {
	var cfg struct {
		Web struct {
			Host            string        `conf:"default:0.0.0.0:8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Qubic struct {
			NodeIps  []string `conf:"default:62.2.219.174"`
			NodePort string   `conf:"default:21841"`
		}
		Opensearch struct {
			Host string `conf:"default:http://93.190.139.223:9200"`
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

	pool := nodes.NewPool(cfg.Qubic.NodeIps)
	if err != nil {
		return errors.Wrap(err, "creating qubic client")
	}

	osClient := opensearch.NewClient(cfg.Opensearch.Host)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      handlers.New(shutdown, log, pool, osClient),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)
	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", cfg.Web.Host)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
