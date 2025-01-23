package main

import (
	"fmt"
	"github.com/pkg/errors"
	qubic "github.com/qubic/go-node-connector"
	rpc "github.com/qubic/qubic-http/foundation/rpc_server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
)

const prefix = "QUBIC_API_SIDECAR"

func main() {
	logger := log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(logger); err != nil {
		logger.Fatalf("main: exited with error: %s", err.Error())
	}
}

func run(logger *log.Logger) error {
	var cfg struct {
		Server struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
			HttpHost        string        `conf:"default:0.0.0.0:8000"`
			GrpcHost        string        `conf:"default:0.0.0.0:8001"`
			MaxTickFetchUrl string        `conf:"default:http://127.0.0.1:8080/max-tick"`
			ReadRetryCount  int           `conf:"default:5"`
		}
		Pool struct {
			NodeFetcherUrl     string        `conf:"default:http://127.0.0.1:8080/status"`
			NodeFetcherTimeout time.Duration `conf:"default:2s"`
			NodePort           string        `conf:"default:21841"`
			InitialCap         int           `conf:"default:5"`
			MaxIdle            int           `conf:"default:20"`
			MaxCap             int           `conf:"default:30"`
			IdleTimeout        time.Duration `conf:"default:15s"`
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
	logger.Printf("main: Config :\n%v\n", out)

	pool, err := qubic.NewPoolConnection(qubic.PoolConfig{
		InitialCap:         cfg.Pool.InitialCap,
		MaxCap:             cfg.Pool.MaxCap,
		MaxIdle:            cfg.Pool.MaxIdle,
		IdleTimeout:        cfg.Pool.IdleTimeout,
		NodeFetcherUrl:     cfg.Pool.NodeFetcherUrl,
		NodeFetcherTimeout: cfg.Pool.NodeFetcherTimeout,
		NodePort:           cfg.Pool.NodePort,
	})
	if err != nil {
		return errors.Wrap(err, "creating qubic pool")
	}

	rpcServer := rpc.NewServer(cfg.Server.GrpcHost, cfg.Server.HttpHost, logger, pool, cfg.Server.MaxTickFetchUrl, cfg.Server.ReadRetryCount)
	err = rpcServer.Start()
	if err != nil {
		return errors.Wrap(err, "starting rpc server")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-shutdown:
			return errors.New("shutting down")
		}
	}
}
