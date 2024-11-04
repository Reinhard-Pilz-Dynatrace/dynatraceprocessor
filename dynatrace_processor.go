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

package dynatraceprocessor

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
		if _, found := attrs.Get(string(MetaDataKeyDTEntityHost)); found {
			continue
		}
		attrs.PutStr(string(MetaDataKeyDTEntityHost), rp.hostID)
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
		if _, found := attrs.Get(string(MetaDataKeyDTEntityHost)); found {
			continue
		}
		attrs.PutStr(string(MetaDataKeyDTEntityHost), rp.hostID)
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
		if _, found := attrs.Get(string(MetaDataKeyDTEntityHost)); found {
			continue
		}
		attrs.PutStr(string(MetaDataKeyDTEntityHost), rp.hostID)
	}
	return ld, nil
}
