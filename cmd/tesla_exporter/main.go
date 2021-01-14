package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uhthomas/tesla_exporter/internal"
	"github.com/uhthomas/tesla_exporter/internal/collector"
	"github.com/uhthomas/tesla_exporter/internal/tesla"
)

func Main(ctx context.Context) error {
	addr := flag.String("addr", ":80", "Listen address.")
	token := flag.String("token", "$TOKEN", "Tesla API token, environment variables are expanded.")
	flag.Parse()

	if *token = os.ExpandEnv(*token); *token == "" {
		return errors.New("token must be set")
	}

	c, err := tesla.New(*token)
	if err != nil {
		return fmt.Errorf("new tesla client: %w", err)
	}

	r := prometheus.NewRegistry()
	if err := r.Register(collector.NewVehicleCollector(ctx, c)); err != nil {
		return fmt.Errorf("register vehicle collector: %w", err)
	}
	return internal.ListenAndServe(ctx, *addr, r)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		select {
		case <-ctx.Done():
		case <-c:
		}
		cancel()
	}()

	if err := Main(ctx); err != nil {
		log.Fatal(err)
	}
}
