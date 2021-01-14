// Copyright 2021, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sumologicexporter

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/consumer/pdata"
	tracetranslator "go.opentelemetry.io/collector/translator/trace"
)

// carbon2TagString  returns all attributes as map of strings
func carbon2TagString(record metricPair) string {
	length := record.attributes.Len()

	if _, exists := record.attributes.Get("metric"); exists {
		length++
	}

	if _, exists := record.attributes.Get("unit"); exists && len(record.metric.Unit()) > 0 {
		length++
	}

	returnValue := make([]string, 0, length)
	record.attributes.ForEach(func(k string, v pdata.AttributeValue) {
		keyName := k

		if k == "name" || k == "unit" {
			keyName = fmt.Sprintf("_%s", k)
		}
		returnValue = append(returnValue, fmt.Sprintf("%s=%s", keyName, tracetranslator.AttributeValueToString(v, false)))
	})

	returnValue = append(returnValue, fmt.Sprintf("metric=%s", record.metric.Name()))

	if len(record.metric.Unit()) > 0 {
		returnValue = append(returnValue, fmt.Sprintf("unit=%s", record.metric.Unit()))
	}

	return strings.Join(returnValue, " ")
}

// carbon2IntRecord converts IntDataPoint to carbon2 metric string
// with additional information from metricPair
func carbon2IntRecord(record metricPair, dataPoint pdata.IntDataPoint) string {
	return fmt.Sprintf("%s  %d %d",
		carbon2TagString(record),
		dataPoint.Value(),
		dataPoint.Timestamp()/1e9,
	)
}

// carbon2DoubleRecord converts DoubleDataPoint to carbon2 metric string
// with additional information from metricPair
func carbon2DoubleRecord(record metricPair, dataPoint pdata.DoubleDataPoint) string {
	return fmt.Sprintf("%s  %g %d",
		carbon2TagString(record),
		dataPoint.Value(),
		dataPoint.Timestamp()/1e9,
	)
}

func metric2Carbon2(record metricPair) string {
	var nextLines []string

	switch record.metric.DataType() {
	case pdata.MetricDataTypeIntGauge:
		for i := 0; i < record.metric.IntGauge().DataPoints().Len(); i++ {
			dataPoint := record.metric.IntGauge().DataPoints().At(i)
			nextLines = append(nextLines, carbon2IntRecord(record, dataPoint))
		}
	case pdata.MetricDataTypeIntSum:
		for i := 0; i < record.metric.IntSum().DataPoints().Len(); i++ {
			dataPoint := record.metric.IntSum().DataPoints().At(i)
			nextLines = append(nextLines, carbon2IntRecord(record, dataPoint))
		}
	case pdata.MetricDataTypeDoubleGauge:
		for i := 0; i < record.metric.DoubleGauge().DataPoints().Len(); i++ {
			dataPoint := record.metric.DoubleGauge().DataPoints().At(i)
			nextLines = append(nextLines, carbon2DoubleRecord(record, dataPoint))
		}
	case pdata.MetricDataTypeDoubleSum:
		for i := 0; i < record.metric.DoubleSum().DataPoints().Len(); i++ {
			dataPoint := record.metric.DoubleSum().DataPoints().At(i)
			nextLines = append(nextLines, carbon2DoubleRecord(record, dataPoint))
		}
	// Skip complex metrics
	case pdata.MetricDataTypeDoubleHistogram:
	case pdata.MetricDataTypeIntHistogram:
	case pdata.MetricDataTypeDoubleSummary:
	}

	return strings.Join(nextLines, "\n")
}
