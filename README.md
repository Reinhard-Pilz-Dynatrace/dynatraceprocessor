# Dynatrace Processor


The Dynatrace processor (config name: dynatrace) adds resource attributes to logs, metrics and traces processed by an OpenTelemetry Collector, so that the Dynatrace can make full use of the ingested data.

## Embedding the Dynatrace Processor into an OpenTelemetry Collector

The Dynatrace Processor is currently not yet included by default in any distribution of the OpenTelemetry Collector.

You need to follow the steps for [building a custom collector](https://opentelemetry.io/docs/collector/custom-collector/). The example below represents a valid `builder-config.yaml` that includes the Dynatrace Processor.

```yaml
dist:
  name: otelcol-dev
  description: Basic OTel Collector distribution that includes the Dynatrace Processor
  output_path: ./otelcol-dev
  otelcol_version: 0.112.0

exporters:
  - gomod: go.opentelemetry.io/collector/exporter/debugexporter v0.112.0
  - gomod: go.opentelemetry.io/collector/exporter/otlpexporter v0.112.0
  - gomod: go.opentelemetry.io/collector/exporter/otlphttpexporter v0.112.0

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.112.0
  - gomod: github.com/Reinhard-Pilz-Dynatrace/dynatraceprocessor v0.112.3

receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.112.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.112.0

providers:
  - gomod: go.opentelemetry.io/collector/confmap/provider/envprovider v1.18.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/fileprovider v1.18.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpprovider v1.18.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.18.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.18.0
```

## Configuration

```yaml
processors:
  dynatrace:
    # Defines whether the `dt.entity.host` resource attribute should be added.
    # default = false
    metadata: {true,false}
```

The example below of a valid `collector-config.yaml` shows how to configure an OpenTelemetry Collector to
* Collect Logs from a file `sample.log`
* Print out the Log Signals on stdout
* Send off the collected logs to Dynatrace
  - The configured `Api-Token` needs to contain the permissions `ingest.logs`

```yaml
receivers:
  filelog:
    include: [ "sample.log" ]
    start_at: beginning

processors:
  dynatrace:
    metadata: true

exporters:
  debug:
    verbosity: detailed
  otlphttp:
    endpoint: "https://########.live.dynatrace.com/api/v2/otlp"
    headers:
      Authorization: "Api-Token dt0c01.########################.################################################################"

service:
  pipelines:
    logs:
      receivers: [filelog]
      processors: [dynatrace]
      exporters: [debug,otlphttp]
```
## Features

### Adding `dt.entity.host` resource attribute
If Dynatrace OneAgent is installed on the host running the OpenTelemetry Collector the resource attribute `dt.entity.host` will be added to the resource attributes of any signal - identifying this specific host as the origin of the OpenTelemetry signals.

Traces, Logs and Metrics already containing the resource attribute `dt.entity.host` will remain untouched.