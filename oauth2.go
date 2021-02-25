package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func newOAuth2Client(ctx context.Context, configPath, tokenPath string) (*http.Client, error) {
	var c oauth2.Config
	if err := openAndDecode(configPath, &c); err != nil {
		return nil, err
	}
	var t oauth2.Token
	if err := openAndDecode(tokenPath, &t); err != nil {
		return nil, err
	}
	return c.Client(ctx, &t), nil
}

func openAndDecode(path string, out interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(out); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}
	return nil
}
