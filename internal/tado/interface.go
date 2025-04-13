package tado

import "context"

type TadoClient interface {
	GetMe(ctx context.Context) (*MeResponse, error)
	GetZones(ctx context.Context, homeId int) (*ZoneResponse, error)
	GetZoneStates(ctx context.Context, homeId int) (*ZoneStateResponse, error)
}
