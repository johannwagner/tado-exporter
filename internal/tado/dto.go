package tado

import "time"

/**
 * Note: Those are DTOs used to parse JSON responses.
 * If this exporter does not use the field, we omit the field altogether to keep
 * the code readable.
 */

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// ME

type MeResponse struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	ID       string  `json:"id"`
	Homes    []Homes `json:"homes"`
	Locale   string  `json:"locale"`
}
type Homes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ZONES
// Note: We only need the name of the zones from here.

type ZoneResponse []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ZONE STATES

type ZoneStateResponse struct {
	ZoneStates ZoneStates `json:"zoneStates"`
}
type ZoneStates map[string]ZoneDetail
type Temperature struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
}
type Setting struct {
	Power       string      `json:"power"`
	Temperature Temperature `json:"temperature"`
}
type HeatingPower struct {
	Percentage float64   `json:"percentage"`
	Timestamp  time.Time `json:"timestamp"`
}
type ActivityDataPoints struct {
	HeatingPower HeatingPower `json:"heatingPower"`
}
type InsideTemperature struct {
	Celsius    float64   `json:"celsius"`
	Fahrenheit float64   `json:"fahrenheit"`
	Timestamp  time.Time `json:"timestamp"`
}
type Humidity struct {
	Percentage float64   `json:"percentage"`
	Timestamp  time.Time `json:"timestamp"`
}
type SensorDataPoints struct {
	InsideTemperature InsideTemperature `json:"insideTemperature"`
	Humidity          Humidity          `json:"humidity"`
}
type ZoneDetail struct {
	TadoMode           string             `json:"tadoMode"`
	Setting            Setting            `json:"setting"`
	ActivityDataPoints ActivityDataPoints `json:"activityDataPoints"`
	SensorDataPoints   SensorDataPoints   `json:"sensorDataPoints"`
}
