package tesla

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	c       *http.Client
	baseURL *url.URL
	token   string
}

func New(token string, opts ...Option) (*Client, error) {
	c := &Client{
		c: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "owner-api.teslamotors.com",
			Path:   "api/v1",
		},
		token: token,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
