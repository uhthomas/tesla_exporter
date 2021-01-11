# Tesla Prometheus Exporter

Export Tesla vehicle metrics to Prometheus.

## Generating a token

Use the [cmd/login](cmd/login) tool.

```sh
bazel run //cmd/login -- --username=hello@tesla.com --passcode=123456
```
