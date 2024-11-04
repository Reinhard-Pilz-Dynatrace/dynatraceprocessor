// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package dynatraceprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/dynatraceprocessor"

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for Resource processor.
type Config struct {
	Metadata bool `mapstructure:"metadata"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
