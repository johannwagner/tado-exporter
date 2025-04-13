package tado

import "context"

type TadoClient interface {
	Authorize() error
	GetMe(ctx context.Context) (*MeResponse, error)
	GetZones(ctx context.Context, homeId int) (*ZoneResponse, error)
	GetZoneStates(ctx context.Context, homeId int) (*ZoneStateResponse, error)
}
