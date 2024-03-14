package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	knauditproxy "github.com/navikt/knaudit-proxy"
)

const (
	envVarOracleURL = "ORACLE_URL"
)

func main() { //nolint: funlen
	var (
		backendType string
		debug       bool
		addr        string
	)

	flag.StringVar(&backendType, "backend-type", "oracle", "Select audit log backend [oracle, stdout]")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.StringVar(&addr, "addr", ":8080", "Address to listen on")
	flag.Parse()

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	slog.Info("starting knaudit-proxy",
		"backend", backendType,
		"addr", addr,
	)

	var (
		backend knauditproxy.SendCloser
		err     error
	)

	switch backendType {
	case "oracle":
		backend = knauditproxy.NewOracleBackend(os.Getenv(envVarOracleURL))
	case "stdout":
		backend = knauditproxy.NewWriterBackend(os.Stdout)
	default:
		slog.Error("unknown backend type", "backend", backendType)
		os.Exit(1)
	}

	err = backend.Open()
	if err != nil {
		slog.Error("opening audit backend", "error", err)
		os.Exit(1)
	}

	defer func() {
		_ = backend.Close()
	}()

	h := knauditproxy.NewServer(backend)

	mux := http.NewServeMux()

	mux.HandleFunc("/", h.StatusHandler)
	mux.HandleFunc("/report", h.ReportHandler)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second, //nolint: gomnd
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("starting server", "error", err.Error())
	}
}
