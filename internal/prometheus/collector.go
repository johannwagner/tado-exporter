package prometheus

import (
	"context"
	"errors"
	"fmt"
	"github.com/johannwagner/tado-exporter-go/internal/tado"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type TadoCollector struct {
	client tado.TadoClient

	tadoClientErrors                   *prometheus.Desc
	tadoActivityHeatingPowerPercentage *prometheus.Desc
	tadoSettingTemperatureValue        *prometheus.Desc
	tadoSensorTemperatureValue         *prometheus.Desc
	tadoSensorHumidityPercentage       *prometheus.Desc
}

func NewTadoCollector(tadoClient tado.TadoClient) *TadoCollector {
	return &TadoCollector{
		client: tadoClient,
		tadoClientErrors: prometheus.NewDesc("tado_client_error",
			"This represents the amount of errors from the tado client",
			[]string{}, nil,
		),
		tadoActivityHeatingPowerPercentage: prometheus.NewDesc("tado_activity_heating_power_percentage",
			"This represent the % of heating power for every zone",
			[]string{"home", "zone"}, nil,
		),
		tadoSettingTemperatureValue: prometheus.NewDesc("tado_setting_temperature_value",
			"This represent the current temperature you asked/programmed in a zone",
			[]string{"home", "zone", "unit"}, nil,
		),
		tadoSensorTemperatureValue: prometheus.NewDesc("tado_sensor_temperature_value",
			"This represent the current temperature detected by sensor in a zone",
			[]string{"home", "zone", "unit"}, nil,
		),

		tadoSensorHumidityPercentage: prometheus.NewDesc("tado_sensor_humidity_percentage",
			"This represent the current humidity % detected by sensor in a zone",
			[]string{"home", "zone"}, nil,
		),
	}
}

func (collector *TadoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.tadoActivityHeatingPowerPercentage
	ch <- collector.tadoSettingTemperatureValue
	ch <- collector.tadoSensorTemperatureValue
	ch <- collector.tadoSensorHumidityPercentage
}

func (collector *TadoCollector) CollectMetrics(ch chan<- prometheus.Metric) error {

	ctx := context.Background()

	meResp, err := collector.client.GetMe(ctx)
	if err != nil {
		return err
	}

	for _, home := range meResp.Homes {
		homeIdStr := strconv.Itoa(home.ID)

		zonesAll, err := collector.client.GetZones(ctx, home.ID)
		if err != nil {
			return err
		}

		zoneStateAll, err := collector.client.GetZoneStates(ctx, home.ID)
		if err != nil {
			return err
		}

		for _, zone := range *zonesAll {
			zoneIdStr := strconv.Itoa(zone.ID)
			zoneState := zoneStateAll.ZoneStates[zoneIdStr]

			ch <- prometheus.MustNewConstMetric(
				collector.tadoActivityHeatingPowerPercentage,
				prometheus.GaugeValue,
				zoneState.ActivityDataPoints.HeatingPower.Percentage,
				homeIdStr,
				zone.Name,
			)

			ch <- prometheus.MustNewConstMetric(
				collector.tadoSettingTemperatureValue,
				prometheus.GaugeValue,
				zoneState.Setting.Temperature.Celsius,
				homeIdStr,
				zone.Name,
				"celsius",
			)

			ch <- prometheus.MustNewConstMetric(
				collector.tadoSettingTemperatureValue,
				prometheus.GaugeValue,
				zoneState.Setting.Temperature.Fahrenheit,
				homeIdStr,
				zone.Name,
				"fahrenheit",
			)

			ch <- prometheus.MustNewConstMetric(
				collector.tadoSensorTemperatureValue,
				prometheus.GaugeValue,
				zoneState.SensorDataPoints.InsideTemperature.Celsius,
				homeIdStr,
				zone.Name,
				"celsius",
			)

			ch <- prometheus.MustNewConstMetric(
				collector.tadoSensorTemperatureValue,
				prometheus.GaugeValue,
				zoneState.SensorDataPoints.InsideTemperature.Fahrenheit,
				homeIdStr,
				zone.Name,
				"fahrenheit",
			)

			ch <- prometheus.MustNewConstMetric(
				collector.tadoSensorHumidityPercentage,
				prometheus.GaugeValue,
				zoneState.SensorDataPoints.Humidity.Percentage,
				homeIdStr,
				zone.Name,
			)

		}
	}
	return nil
}

// Collect implements required collect function for all promehteus collectors
func (collector *TadoCollector) Collect(ch chan<- prometheus.Metric) {

	err := collector.CollectMetrics(ch)
	clientError := 0
	if err != nil {
		if errors.Is(err, tado.ClientNotInitializedError{}) {
			clientError = 1
		}
		fmt.Println(err)
	}

	ch <- prometheus.MustNewConstMetric(
		collector.tadoClientErrors,
		prometheus.GaugeValue,
		float64(clientError),
	)
}
