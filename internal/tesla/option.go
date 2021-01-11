package tesla

import (
	"net/http"
	urlpkg "net/url"
)

type Option func(*Client) error

func HTTPClient(c *http.Client) Option {
	return func(cc *Client) (err error) {
		cc.c = c
		return nil
	}
}

func BaseURL(url string) Option {
	return func(c *Client) (err error) {
		c.baseURL, err = urlpkg.Parse(url)
		return err
	}
}
