// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xsentinel

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/core/log"
	"github.com/alibaba/sentinel-golang/core/stat"
	"github.com/alibaba/sentinel-golang/core/system"
)

var globalSlotChain = BuildDefaultSlotChain()

func Init(opt *Option) error {
	// TODO: 考虑放入配置文件
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = NewLogger()
	// 关闭监控日志文件异步输出
	conf.Sentinel.Log.Metric.FlushIntervalSec = 0
	err := api.InitWithConfig(conf)
	if err != nil {
		return err
	}

	// circuit breaker
	if len(opt.CircuitBreakerRules) > 0 {
		_, err = circuitbreaker.LoadRules(opt.CircuitBreakerRules)
		if err != nil {
			return err
		}
	}

	// flow
	if len(opt.FlowRules) > 0 {
		_, err = flow.LoadRules(opt.FlowRules)
		if err != nil {
			return err
		}
	}

	// hotspot
	if len(opt.HotspotRules) > 0 {
		_, err = hotspot.LoadRules(opt.HotspotRules)
		if err != nil {
			return err
		}
	}

	// isolation
	if len(opt.IsolationRules) > 0 {
		_, err = isolation.LoadRules(opt.IsolationRules)
		if err != nil {
			return err
		}
	}

	// system
	if len(opt.SystemRules) > 0 {
		_, err = system.LoadRules(opt.SystemRules)
		if err != nil {
			return err
		}
	}

	return nil
}

func Entry(resource string, opts ...api.EntryOption) (*base.SentinelEntry, *base.BlockError) {
	opts = append(opts, api.WithSlotChain(GlobalSlotChain()))
	return api.Entry(resource, opts...)
}

func TraceError(entry *base.SentinelEntry, err error) {
	api.TraceError(entry, err)
}

func GlobalSlotChain() *base.SlotChain {
	return globalSlotChain
}

func BuildDefaultSlotChain() *base.SlotChain {
	sc := base.NewSlotChain()
	sc.AddStatPrepareSlot(stat.DefaultResourceNodePrepareSlot)

	sc.AddRuleCheckSlot(system.DefaultAdaptiveSlot)
	sc.AddRuleCheckSlot(flow.DefaultSlot)
	sc.AddRuleCheckSlot(isolation.DefaultSlot)
	sc.AddRuleCheckSlot(hotspot.DefaultSlot)
	sc.AddRuleCheckSlot(circuitbreaker.DefaultSlot)

	sc.AddStatSlot(stat.DefaultSlot)
	sc.AddStatSlot(log.DefaultSlot)
	sc.AddStatSlot(flow.DefaultStandaloneStatSlot)
	sc.AddStatSlot(hotspot.DefaultConcurrencyStatSlot)
	sc.AddStatSlot(circuitbreaker.DefaultMetricStatSlot)

	// nss metrics
	// sc.AddStatPrepareSlot(NssMetricSlot)
	// sc.AddStatSlot(NssMetricSlot)
	return sc
}
