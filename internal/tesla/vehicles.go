package tesla

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type vehiclesResponse struct {
	Response []struct {
		ID uint64 `json:"id"`
	} `json:"response"`
	Count int `json:"count"`
}

// Vehicles lists all vehicles associated with the account, and describes them
// in detail.
func (c *Client) Vehicles(ctx context.Context) ([]*Vehicle, error) {
	u := *c.baseURL
	u.Path = path.Join(u.Path, "vehicles")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "tesla_exporter")

	res, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var out vehiclesResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	vs := make([]*Vehicle, len(out.Response))
	for i, r := range out.Response {
		v, err := c.Vehicle(ctx, r.ID)
		if err != nil {
			return nil, fmt.Errorf("get vehicle %d: %w", r.ID, err)
		}
		vs[i] = v
	}
	return vs, nil
}
