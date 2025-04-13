package tado

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

type TadoAPIClient struct {
	client *http.Client
}

func NewTadoAPIClient() *TadoAPIClient {

	return &TadoAPIClient{
		client: nil,
	}
}

func (t *TadoAPIClient) Authorize() error {
	// Taken from https://support.tado.com/en/articles/8565472-how-do-i-authenticate-to-access-the-rest-api
	config := oauth2.Config{
		ClientID: DEVICE_FLOW_CLIENT_ID,
		Scopes:   []string{"offline_access"},

		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: DEVICE_AUTH_URL,
			TokenURL:      TOKEN_URL,
		},
	}
	ctx := context.Background()

	deviceCode, err := config.DeviceAuth(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Go to %v and enter code %v\n", deviceCode.VerificationURIComplete, deviceCode.UserCode)

	token, err := config.DeviceAccessToken(
		ctx,
		deviceCode,
	)
	if err != nil {
		return err
	}

	fmt.Printf("Login successful, starting...\n")
	t.client = config.Client(ctx, token)
	return nil
}

func (t *TadoAPIClient) GetJSON(ctx context.Context, url string, res interface{}) error {
	fullUrl := fmt.Sprintf("%v%v", API_URL, url)
	req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)

	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "tado-exporter-go - https://github.com/johannwagner/tado-exporter-go")
	r, err := t.client.Do(req)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(res)
}

func (t *TadoAPIClient) GetMe(ctx context.Context) (*MeResponse, error) {
	var resp MeResponse
	err := t.GetJSON(ctx, "/me", &resp)
	return &resp, err
}

func (t *TadoAPIClient) GetZones(ctx context.Context, homeId int) (*ZoneResponse, error) {
	var resp ZoneResponse
	err := t.GetJSON(ctx, "/homes/"+strconv.Itoa(homeId)+"/zones", &resp)
	return &resp, err
}

func (t *TadoAPIClient) GetZoneStates(ctx context.Context, homeId int) (*ZoneStateResponse, error) {
	var resp ZoneStateResponse
	err := t.GetJSON(ctx, "/homes/"+strconv.Itoa(homeId)+"/zoneStates", &resp)
	return &resp, err
}
