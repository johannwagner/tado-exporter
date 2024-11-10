package tado

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"strconv"
	"time"
)

type TadoClient struct {
	tadoUsername string
	tadoPassword string

	accessToken       *string
	accessTokenExpiry *time.Time

	refreshToken *string
}

var LEAKED_CLIENT_SECRET = "wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc"

var AUTH_URL = "https://auth.tado.com/oauth/token"
var API_URL = "https://my.tado.com/api/v2"

func NewTadoClient(tadoUsername string, tadoPassword string) *TadoClient {

	return &TadoClient{
		tadoUsername:      tadoUsername,
		tadoPassword:      tadoPassword,
		accessToken:       nil,
		accessTokenExpiry: nil,
		refreshToken:      nil,
	}

}

func (tc *TadoClient) processAuthResponse(ar *AuthResponse) {

	fmt.Println(ar.AccessToken)

	tc.accessToken = &ar.AccessToken
	tc.refreshToken = &ar.RefreshToken

	expCalc := time.Now().Add(time.Second * time.Duration(ar.ExpiresIn))
	tc.accessTokenExpiry = &expCalc
}

func (tc *TadoClient) RefreshToken(ctx context.Context) error {
	var output AuthResponse
	err := requests.URL(AUTH_URL).
		BodyForm(map[string][]string{
			"client_id":     {"tado-web-app"},
			"client_secret": {LEAKED_CLIENT_SECRET},
			"grant_type":    {"refresh_token"},
			"scope":         {"home.user"},
			"refresh_token": {*tc.refreshToken},
		}).ToJSON(&output).Fetch(ctx)

	if err != nil {
		return err
	}

	tc.processAuthResponse(&output)

	return nil
}

func (tc *TadoClient) AskForToken(ctx context.Context) error {

	var output AuthResponse

	err := requests.URL(AUTH_URL).
		BodyForm(map[string][]string{
			"client_id":     {"tado-web-app"},
			"client_secret": {LEAKED_CLIENT_SECRET},
			"grant_type":    {"password"},
			"scope":         {"home.user"},
			"username":      {tc.tadoUsername},
			"password":      {tc.tadoPassword},
		}).ToJSON(&output).Fetch(ctx)

	if err != nil {
		return err
	}

	tc.processAuthResponse(&output)

	return nil
}

func (tc *TadoClient) EnsureAuthentication(ctx context.Context) error {

	// We have a valid token, go ahead.
	if tc.accessTokenExpiry != nil && (*tc.accessTokenExpiry).After(time.Now()) {
		return nil
	}

	if tc.refreshToken != nil {
		return tc.RefreshToken(ctx)
	}

	return tc.AskForToken(ctx)
}

func (tc *TadoClient) GetMe(ctx context.Context) (*MeResponse, error) {
	err := tc.EnsureAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	var resp MeResponse
	err = requests.URL(API_URL + "/me").Bearer(*tc.accessToken).ToJSON(&resp).Fetch(ctx)
	return &resp, err
}

func (tc *TadoClient) GetZones(ctx context.Context, homeId int) (*ZoneResponse, error) {
	err := tc.EnsureAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	var resp ZoneResponse
	err = requests.URL(API_URL + "/homes/" + strconv.Itoa(homeId) + "/zones").Bearer(*tc.accessToken).ToJSON(&resp).Fetch(ctx)
	return &resp, err
}

func (tc *TadoClient) GetZoneStates(ctx context.Context, homeId int) (*ZoneStateResponse, error) {
	err := tc.EnsureAuthentication(ctx)
	if err != nil {
		return nil, err
	}
	var resp ZoneStateResponse
	err = requests.URL(API_URL + "/homes/" + strconv.Itoa(homeId) + "/zoneStates").Bearer(*tc.accessToken).ToJSON(&resp).Fetch(ctx)
	return &resp, err
}
