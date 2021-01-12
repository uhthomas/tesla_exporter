package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	c                       *http.Client
	challenge, challengeSum string
}

func newClient() (*Client, error) {
	var b [86]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		return nil, fmt.Errorf("rand read full: %w", err)
	}
	challenge := base64.RawURLEncoding.EncodeToString(b[:])
	sum := sha256.Sum256([]byte(challenge))

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("new cookie jar: %w", err)
	}

	return &Client{
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
			Jar: jar,
		},
		challenge:    challenge,
		challengeSum: base64.RawURLEncoding.EncodeToString(sum[:]),
	}, nil
}

func (c *Client) Login(ctx context.Context, username, password, passcode string) (token string, err error) {
	fmt.Println("authenticating")
	transactionID, err := c.authenticate(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("authenticate: %w", err)
	}

	fmt.Println("listing devices")
	devices, err := c.listDevices(ctx, transactionID)
	if err != nil {
		return "", fmt.Errorf("list devices: %w", err)
	}

	if len(devices) == 0 {
		return "", errors.New("no devices")
	}

	for _, d := range devices {
		fmt.Printf("verifying device: %s\n", d.Name)
		if err := c.verify(ctx, transactionID, d.ID, passcode); err != nil {
			return "", fmt.Errorf("verify: %w", err)
		}
	}

	fmt.Println("authorizing")
	code, err := c.authorize(ctx, transactionID)
	if err != nil {
		return "", fmt.Errorf("authorize: %w", err)
	}

	fmt.Println("access token")
	token, err = c.accessToken(ctx, code)
	if err != nil {
		return "", fmt.Errorf("access token: %w", err)
	}

	fmt.Println("exchanging token")
	return c.exchangeToken(ctx, token)
}

func (c *Client) authURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "auth.tesla.com",
		Path:   "oauth2/v3/authorize",
		RawQuery: url.Values{
			"client_id":             {"ownerapi"},
			"audience":              {""},
			"code_challenge":        {c.challengeSum},
			"code_challenge_method": {"S256"},
			"locale":                {"en"},
			"prompt":                {"login"},
			"redirect_uri":          {"https://auth.tesla.com/void/callback"},
			"response_type":         {"code"},
			"scope":                 {"openid email offline_access"},
			"state":                 {"tesla_exporter"},
		}.Encode(),
	}
}

func (c *Client) authenticate(ctx context.Context, username, password string) (transactionID string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.authURL().String(), nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "tesla_exporter")

	res, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	d, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", fmt.Errorf("new document: %w", err)
	}

	v := url.Values{
		"identity":   {username},
		"credential": {password},
	}

	d.Find("input[type=hidden]").Each(func(_ int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if !ok {
			return
		}
		value, ok := s.Attr("value")
		if !ok {
			return
		}
		v.Set(name, value)
	})

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, c.authURL().String(), strings.NewReader(v.Encode()))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "tesla_exporter")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return v.Get("transaction_id"), nil
}

type Device struct {
	DispatchRequired bool      `json:"dispatchRequired"`
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	FactorType       string    `json:"factorType"`
	FactorProvider   string    `json:"factorProvider"`
	SecurityLevel    int       `json:"securityLevel"`
	Activated        time.Time `json:"activatedAt"`
	Updated          time.Time `json:"updatedAt"`
}

func (c *Client) listDevices(ctx context.Context, transactionID string) ([]Device, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, (&url.URL{
		Scheme:   "https",
		Host:     "auth.tesla.com",
		Path:     "oauth2/v3/authorize/mfa/factors",
		RawQuery: url.Values{"transaction_id": {transactionID}}.Encode(),
	}).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "tesla_exporter")

	res, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var out struct {
		Data []Device `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return out.Data, nil
}

func (c *Client) verify(ctx context.Context, transactionID, deviceID, passcode string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]string{
		"transaction_id": transactionID,
		"factor_id":      deviceID,
		"passcode":       passcode,
	}); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://auth.tesla.com/oauth2/v3/authorize/mfa/verify", &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "tesla_exporter")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	var out struct {
		Data struct {
			Approved bool `json:"approved"`
		} `json:"data`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	if !out.Data.Approved {
		return errors.New("not approved")
	}
	return nil
}

func (c *Client) authorize(ctx context.Context, transactionID string) (code string, err error) {
	cr := c.c.CheckRedirect
	c.c.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() { c.c.CheckRedirect = cr }()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.authURL().String(), strings.NewReader(url.Values{
		"transaction_id": {transactionID},
	}.Encode()))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "tesla_exporter")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	u, err := res.Location()
	if err != nil {
		return "", fmt.Errorf("response location: %w", err)
	}
	return u.Query().Get("code"), nil
}

func (c *Client) accessToken(ctx context.Context, code string) (token string, err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     "ownerapi",
		"code_verifier": c.challenge,
		"code":          code,
		"redirect_uri":  "https://auth.tesla.com/void/callback",
	}); err != nil {
		return "", fmt.Errorf("json encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://auth.tesla.com/oauth2/v3/token", &buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "tesla_exporter")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("json decode: %w", err)
	}
	return out.AccessToken, nil
}

func (c *Client) exchangeToken(ctx context.Context, token string) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]string{
		"grant_type":    "urn:ietf:params:oauth:grant-type:jwt-bearer",
		"client_id":     "81527cff06843c8634fdc09e8ac0abefb46ac849f38fe1e431c2ef2106796384",
		"client_secret": "c7257eb71a564034f9419ee651c7d0e5f7aa6bfbd18bafb5c5c033b093bb2fa3",
	}); err != nil {
		return "", fmt.Errorf("json encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://owner-api.teslamotors.com/oauth/token", &buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "tesla_exporter")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("json decode: %w", err)
	}
	return out.AccessToken, nil
}
