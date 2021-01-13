package tesla

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type vehiclesResponse struct {
	Response Vehicles `json:"response"`
	Count    int      `json:"count"`
}

type Vehicles []struct {
	ID                     uint64      `json:"id"`
	VehicleID              uint64      `json:"vehicle_id"`
	VIN                    string      `json:"vin"`
	DisplayName            string      `json:"display_name"`
	OptionCodes            string      `json:"option_codes"`
	Color                  interface{} `json:"color"`
	AccessType             string      `json:"access_type"`
	Tokens                 []string    `json:"tokens"`
	State                  string      `json:"state"`
	InService              bool        `json:"in_service"`
	IDS                    string      `json:"id_s"`
	CalendarEnabled        bool        `json:"calendar_enabled"`
	APIVersion             int         `json:"api_version"`
	BackseatToken          interface{} `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{} `json:"backseat_token_updated_at"`
	VehicleConfig          interface{} `json:"vehicle_config"`
}

// Vehicles lists all vehicles associated with the account, and describes them
// in detail.
func (c *Client) Vehicles(ctx context.Context) (Vehicles, error) {
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
	return out.Response, nil
}
