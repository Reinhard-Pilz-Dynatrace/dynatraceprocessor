# Dynatrace Processor

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [beta]: traces, metrics, logs   |
| Distributions | [core], [contrib], [k8s] |
| [Code Owners](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/CONTRIBUTING.md#becoming-a-code-owner)    | [@Reinhard-Pilz-Dynatrace](https://github.com/Reinhard-Pilz-Dynatrace) |

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[core]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[k8s]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-k8s
<!-- end autogenerated section -->

The Dynatrace processor (config name: dynatrace) adds resource attributes to logs, metrics and traces sent to Dynatrace so that the Dynatrace can make full use of the ingested data.

## Configuration

```yaml
processors:
  dynatrace:
    # Defines whether the `dt.entity.host` resource attribute should be added.
    # default = false
    metadata: {true,false}
```
## Features

### Adding `dt.entity.host` resource attribute
If Dynatrace OneAgent is installed on the host running the OpenTelemetry Collector the resource attribute `dt.entity.host` will be added to the resource attributes of any signal - identifying this specific host as the origin of the OpenTelemetry signals.

Traces, Logs and Metrics already containing the resource attribute `dt.entity.host` will remain untouched.