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
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"

	telegraf "github.com/pmalek-sumo/telegraf/agent"
	"go.uber.org/zap"
)

type telegrafreceiver struct {
	sync.Mutex
	startOnce sync.Once
	stopOnce  sync.Once
	wg        sync.WaitGroup
	cancel    context.CancelFunc

	// agent  *stanza.LogAgent
	agent *telegraf.Agent
	// emitter  *LogEmitter
	consumer consumer.MetricsConsumer
	logger   *zap.Logger
}

// Ensure this receiver adheres to required interface
var _ component.MetricsReceiver = (*telegrafreceiver)(nil)

// Start tells the receiver to start
func (r *telegrafreceiver) Start(ctx context.Context, host component.Host) error {
	r.Lock()
	defer r.Unlock()
	err := componenterror.ErrAlreadyStarted
	r.startOnce.Do(func() {
		err = nil
		rctx, cancel := context.WithCancel(ctx)
		r.cancel = cancel
		r.logger.Info("Starting telegraf receiver")

		// if obsErr := r.agent.Run(rctx); obsErr != nil {
		// 	err = fmt.Errorf("start telegaf receiver: %s", obsErr)
		// 	return
		// }

		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case <-rctx.Done():
					return
				case <-time.After(time.Second):
					r.logger.Info("ping")
					// case obsLog, ok := <-r.emitter.metricsChan:
					// 	if !ok {
					// 		continue
					// 	}

					// 	if consumeErr := r.consumer.ConsumeMetrics(ctx, nil); consumeErr != nil {
					// 		r.logger.Error("ConsumeLogs() error", zap.String("error", consumeErr.Error()))
					// 	}
				}
			}
		}()
	})

	return err
}

// Shutdown is invoked during service shutdown
func (r *telegrafreceiver) Shutdown(context.Context) error {
	r.Lock()
	defer r.Unlock()

	err := componenterror.ErrAlreadyStopped
	r.stopOnce.Do(func() {
		r.logger.Info("Stopping telegraf receiver")
		r.cancel()
		r.wg.Wait()
		err = nil
	})
	return err
}
