# tado-exporter 

This is a rewrite based on `eko/tado-exporter`, because it stopped working
after a tado API change in the beginning of November 2024.
Since I am better and faster in writing Go, I did it in Go and reimplemented the existing API.

This should be a drop-in for the most metrics.

## Available environment variables

| Environment variable name    | Description                                                                                |
|:----------------------------:|--------------------------------------------------------------------------------------------|
| EXPORTER_USERNAME      | Required. This represent your tado° account username/email                                       |
| EXPORTER_PASSWORD      | Required. This represent your tado° account password                                             |

## Available Prometheus metrics

| Metric name                  | Description                                                                                |
|:----------------------------:|--------------------------------------------------------------------------------------------|
| tado_activity_heating_power_percentage | This represent the % of heating power for every zone                             |
| tado_setting_temperature_value         | This represent the current temperature you asked/programmed in a zone            |
| tado_sensor_temperature_value          | This represent the current temperature detected by sensor in a zone              |
| tado_sensor_humidity_percentage        | This represent the current humidity % detected by sensor in a zone               |