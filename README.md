# Tesla Prometheus Exporter

Export Tesla vehicle metrics to Prometheus.

## Getting started

To get started, [generate a token](#generated-a-token). Container images are
pushed to the GitHub container registry, and tagged by release version. To
use the latest version, please see the [releases](releases).

A token must be provided either as an environment variable, or as a flag.

```sh
docker run -e TOKEN=abc ghcr.io/uhthomas/tesla_exporter
```

## Generating a token

Use the [cmd/login](cmd/login) tool.

```sh
bazel run //cmd/login -- --username=hello@tesla.com --passcode=123456
```

You'll then be prompted to enter your password, after which you'll be presented
with your token.

## Metrics

All metrics are labeled using the car's VIN. These are currently non-exhaustive,
as many more are planned to be added.

| Metric                                | API reference                       |
| :------------------------------------ | :---------------------------------- |
| tesla_vehicle_info                    | id, vehicle_id                      |
| tesla_vehicle_name                    | display_name                        |
| tesla_vehicle_state                   | state                               |
| tesla_vehicle_software_version        | vehicle_state.car_version           |
| tesla_vehicle_odometer_miles_total    | vehicle_state.odometer              |
| tesla_vehicle_inside_temp_celsius     | climate_state.inside_temp           |
| tesla_vehicle_outside_temp_celsius    | climate_state.outside_temp          |
| tesla_vehicle_battery_ratio           | charge_state.battery_level          |
| tesla_vehicle_battery_usable_ratio    | charge_state.usable_battery_level   |
| tesla_vehicle_battery_ideal_miles     | charge_state.battery_range          |
| tesla_vehicle_battery_estimated_miles | charge_state.est_battery_range      |
| tesla_vehicle_charge_volts            | charge_state.charger_voltage        |
| tesla_vehicle_charge_amps             | charge_state.charger_actual_current |
| tesla_vehicle_charge_amps_available   | charge_state.charger_pilot_current  |

## Preview

The metrics collected by Prometheus can be visualized through Grafana, and can
look something like the following.

![Grafana](docs/images/grafana.png)
