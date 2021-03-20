package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uhthomas/tesla"
	"github.com/uhthomas/tesla_exporter/internal"
	"github.com/uhthomas/tesla_exporter/internal/collector"
)

func Main(ctx context.Context) error {
	addr := flag.String("addr", ":80", "Listen address.")
	oauth2ConfigPath := flag.String("oauth2-config-path", "oauth2_config.json", "Tesla OAuth2 config file")
	oauth2TokenPath := flag.String("oauth2-token-path", "oauth2_token.json", "Tesla OAuth2 token file")
	expire := flag.Duration("expire", 30*time.Second, "Expire cache metrics")
	flag.Parse()

	s, ctx, err := tesla.New(ctx, &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	},
		tesla.OAuth2(*oauth2ConfigPath, *oauth2TokenPath),
	)
	if err != nil {
		return fmt.Errorf("new tesla service: %w", err)
	}

	r := prometheus.NewRegistry()
	if err := r.Register(collector.NewVehicleCollector(ctx, s, *expire)); err != nil {
		return fmt.Errorf("register vehicle collector: %w", err)
	}
	return internal.ListenAndServe(ctx, *addr, r)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := Main(ctx); err != nil {
		log.Fatal(err)
	}
}
