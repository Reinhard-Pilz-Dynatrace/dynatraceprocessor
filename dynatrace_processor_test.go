/**
 * @license
 * Copyright 2020 Dynatrace LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dynatraceprocessor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor/processortest"

	"github.com/Reinhard-Pilz-Dynatrace/dynatraceprocessor"
	"github.com/Reinhard-Pilz-Dynatrace/dynatraceprocessor/testdata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/ptracetest"
)

func TestDynatraceProcessorAttributesInsert(t *testing.T) {
	const mockEvalDTEntityHost = "HOST-2EF98EFF909EE3F6"
	const mockConfDTEntityHost = "HOST-0000000000000000"
	tests := []struct {
		name             string
		config           *dynatraceprocessor.Config
		sourceAttributes map[string]string
		wantAttributes   map[string]string
	}{
		{
			name:             "config_with_attribute_applied_on_nil_resource_enabled",
			config:           &dynatraceprocessor.Config{Metadata: true},
			sourceAttributes: nil,
			wantAttributes:   map[string]string{dynatraceprocessor.KeyEntityHost: mockEvalDTEntityHost},
		},
		{
			name:             "config_with_attribute_applied_on_empty_resource_enabled",
			config:           &dynatraceprocessor.Config{Metadata: true},
			sourceAttributes: map[string]string{},
			wantAttributes:   map[string]string{dynatraceprocessor.KeyEntityHost: mockEvalDTEntityHost},
		},
		{
			name:             "config_attribute_applied_on_existing_resource_attributes_enabled",
			config:           &dynatraceprocessor.Config{Metadata: true},
			sourceAttributes: map[string]string{dynatraceprocessor.KeyEntityHost: mockConfDTEntityHost},
			wantAttributes:   map[string]string{dynatraceprocessor.KeyEntityHost: mockConfDTEntityHost},
		},
		{
			name:             "config_with_attribute_applied_on_nil_resource_disabled",
			config:           &dynatraceprocessor.Config{Metadata: false},
			sourceAttributes: nil,
			wantAttributes:   map[string]string{},
		},
		{
			name:             "config_with_attribute_applied_on_empty_resource_disabled",
			config:           &dynatraceprocessor.Config{Metadata: false},
			sourceAttributes: map[string]string{},
			wantAttributes:   map[string]string{},
		},
		{
			name:             "config_attribute_applied_on_existing_resource_attributes_disabled",
			config:           &dynatraceprocessor.Config{Metadata: false},
			sourceAttributes: map[string]string{dynatraceprocessor.KeyEntityHost: mockConfDTEntityHost},
			wantAttributes:   map[string]string{dynatraceprocessor.KeyEntityHost: mockConfDTEntityHost},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test trace consumer
			ttn := new(consumertest.TracesSink)

			ctx := context.WithValue(context.Background(), dynatraceprocessor.MetaDataKeyDTEntityHost, mockEvalDTEntityHost)

			factory := dynatraceprocessor.NewFactory()
			rtp, err := factory.CreateTraces(ctx, processortest.NewNopSettings(), tt.config, ttn)
			require.NoError(t, err)
			assert.True(t, rtp.Capabilities().MutatesData)

			sourceTraceData := generateTraceData(tt.sourceAttributes)
			wantTraceData := generateTraceData(tt.wantAttributes)
			err = rtp.ConsumeTraces(ctx, sourceTraceData)
			require.NoError(t, err)
			traces := ttn.AllTraces()
			require.Len(t, traces, 1)
			assert.NoError(t, ptracetest.CompareTraces(wantTraceData, traces[0]))

			// Test metrics consumer
			tmn := new(consumertest.MetricsSink)
			rmp, err := factory.CreateMetrics(ctx, processortest.NewNopSettings(), tt.config, tmn)
			require.NoError(t, err)
			assert.True(t, rtp.Capabilities().MutatesData)

			sourceMetricData := generateMetricData(tt.sourceAttributes)
			wantMetricData := generateMetricData(tt.wantAttributes)
			err = rmp.ConsumeMetrics(ctx, sourceMetricData)
			require.NoError(t, err)
			metrics := tmn.AllMetrics()
			require.Len(t, metrics, 1)
			assert.NoError(t, pmetrictest.CompareMetrics(wantMetricData, metrics[0]))

			// Test logs consumer
			tln := new(consumertest.LogsSink)
			rlp, err := factory.CreateLogs(ctx, processortest.NewNopSettings(), tt.config, tln)
			require.NoError(t, err)
			assert.True(t, rtp.Capabilities().MutatesData)

			sourceLogData := generateLogData(tt.sourceAttributes)
			wantLogData := generateLogData(tt.wantAttributes)
			err = rlp.ConsumeLogs(ctx, sourceLogData)
			require.NoError(t, err)
			logs := tln.AllLogs()
			require.Len(t, logs, 1)
			assert.NoError(t, plogtest.CompareLogs(wantLogData, logs[0]))
		})
	}
}

func generateTraceData(attributes map[string]string) ptrace.Traces {
	td := testdata.GenerateTracesOneSpanNoResource()
	if attributes == nil {
		return td
	}
	resource := td.ResourceSpans().At(0).Resource()
	for k, v := range attributes {
		resource.Attributes().PutStr(k, v)
	}
	return td
}

func generateMetricData(attributes map[string]string) pmetric.Metrics {
	md := testdata.GenerateMetricsOneMetricNoResource()
	if attributes == nil {
		return md
	}
	resource := md.ResourceMetrics().At(0).Resource()
	for k, v := range attributes {
		resource.Attributes().PutStr(k, v)
	}
	return md
}

func generateLogData(attributes map[string]string) plog.Logs {
	ld := testdata.GenerateLogsOneLogRecordNoResource()
	if attributes == nil {
		return ld
	}
	resource := ld.ResourceLogs().At(0).Resource()
	for k, v := range attributes {
		resource.Attributes().PutStr(k, v)
	}
	return ld
}
