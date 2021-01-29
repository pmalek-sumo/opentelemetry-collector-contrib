// Copyright 2019, OpenTelemetry Authors
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

package telegrafreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"

	telegraf "github.com/pmalek-sumo/telegraf/agent"
)

const (
	typeStr = "telegraf"
	verStr  = "0.12.0"
)

// NewFactory creates a factory for Stanza receiver
// TODO
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithMetrics(createMetricsReceiver),
	)
}

func createDefaultConfig() configmodels.Receiver {
	return nil
}

// CreateLogsReceiver creates a logs receiver based on provided config
// TODO
func createMetricsReceiver(
	ctx context.Context,
	params component.ReceiverCreateParams,
	cfg configmodels.Receiver,
	nextConsumer consumer.MetricsConsumer,
) (component.MetricsReceiver, error) {

	// obsConfig := cfg.(*Config)

	// emitter := NewLogEmitter(params.Logger.Sugar())

	// pipeline, err := obsConfig.Operators.IntoPipelineConfig()
	// if err != nil {
	// 	return nil, err
	// }
	var err error

	agent, err := telegraf.NewAgent(nil)
	_ = agent
	// logAgent, err := stanza.NewBuilder(params.Logger.Sugar()).
	// 	WithConfig(&stanza.Config{Pipeline: pipeline}).
	// 	WithPluginDir(obsConfig.PluginDir).
	// 	WithDatabaseFile(obsConfig.OffsetsFile).
	// 	WithDefaultOutput(emitter).
	// 	Build()
	if err != nil {
		return nil, err
	}

	return &telegrafreceiver{
		agent: agent,
		// 	emitter:  emitter,
		// 	consumer: nextConsumer,
		// 	logger:   params.Logger,
	}, nil
}
