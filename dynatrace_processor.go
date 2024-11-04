// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package dynatraceprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/dynatraceprocessor"

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type dynatraceProcessor struct {
	logger *zap.Logger
	hostID string
}

func (rp *dynatraceProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if len(rp.hostID) == 0 {
		return td, nil
	}
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		attrs := rss.At(i).Resource().Attributes()
		if _, found := attrs.Get(string(metaDataKeyDTEntityHost)); found {
			continue
		}
		attrs.PutStr(string(metaDataKeyDTEntityHost), rp.hostID)
	}
	return td, nil
}

func (rp *dynatraceProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	if len(rp.hostID) == 0 {
		return md, nil
	}
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		attrs := rms.At(i).Resource().Attributes()
		if _, found := attrs.Get(string(metaDataKeyDTEntityHost)); found {
			continue
		}
		attrs.PutStr(string(metaDataKeyDTEntityHost), rp.hostID)
	}
	return md, nil
}

func (rp *dynatraceProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	if len(rp.hostID) == 0 {
		return ld, nil
	}
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		attrs := rls.At(i).Resource().Attributes()
		if _, found := attrs.Get(string(metaDataKeyDTEntityHost)); found {
			continue
		}
		attrs.PutStr(string(metaDataKeyDTEntityHost), rp.hostID)
	}
	return ld, nil
}
