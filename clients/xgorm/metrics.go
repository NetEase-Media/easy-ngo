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

package xgorm

import (
	"strings"
	"time"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

// 数据库常用的指标
// 参考维度
// sql id 和 查询语句
// sql id | sqlString | 调用次数 | 总时间 | 平均RT | 最大并发 | 最慢调用 | 错误次数 | 影响行数 | 读取行数 | 平均一次读取行数

// ms0_10 | ms10_100 | ms100_1000 | s1_10 | s10_n

const (
	ngoMetricsKey = "ngo:db:metrics"
)

const (
	metricRequestTotalName         = "gorm_request_total"
	metricRequestDurationName      = "gorm_request_duration"
	metricRequestDurationRangeName = "gorm_request_duration_range"
	metricRequestDurationAllName   = "gorm_request_duration_all"
	metricRequestErrorName         = "gorm_request_error"
	metricRequestEffectRowName     = "gorm_request_effect_row"
	metricRequestReadRowName       = "gorm_request_read_row"
)

var (
	metricRequestTotal         metrics.Counter
	metricRequestDuration      metrics.Gauge
	metricRequestDurationRange metrics.Histogram
	metricRequestDurationAll   metrics.Counter
	metricRequestError         metrics.Counter
	metricRequestEffectRow     metrics.Gauge
	metricRequestReadRow       metrics.Gauge
	labelValues                = []string{"dsn", "sqlString"}
)

type gormMetricsPlugin struct {
	enable  bool
	metrics metrics.Provider
}

func newGormMetricsPlugin(enable bool, metrics metrics.Provider) *gormMetricsPlugin {
	if enable {
		metricRequestTotal = metrics.NewCounter(metricRequestTotalName, labelValues...)
		metricRequestDuration = metrics.NewGauge(metricRequestDurationName, labelValues...)
		bukets := prometheus.ExponentialBuckets(.01, 10, 5)
		metricRequestDurationRange = metrics.NewHistogram(metricRequestDurationRangeName, bukets, labelValues...)
		metricRequestDurationAll = metrics.NewCounter(metricRequestDurationAllName, labelValues...)
		metricRequestError = metrics.NewCounter(metricRequestErrorName, labelValues...)
		metricRequestEffectRow = metrics.NewGauge(metricRequestEffectRowName, labelValues...)
		metricRequestReadRow = metrics.NewGauge(metricRequestReadRowName, labelValues...)
	}
	return &gormMetricsPlugin{
		enable:  enable,
		metrics: metrics,
	}
}

func (p *gormMetricsPlugin) Name() string {
	return ngoMetricsKey
}

func (p *gormMetricsPlugin) Initialize(db *gorm.DB) error {
	p.registerCallbacks(db)
	return nil
}

func (p *gormMetricsPlugin) registerCallbacks(db *gorm.DB) {

	db.Callback().Query().Before("gorm:query").Register("ngo:metrics:before_query", p.metricBefore)
	db.Callback().Query().After("gorm:query").Register("ngo:metrics:after_query", p.metricAfter)

	db.Callback().Create().Before("gorm:create").Register("ngo:metrics:before_create", p.metricBefore)
	db.Callback().Create().After("gorm:create").Register("ngo:metrics:after_create", p.metricAfter)

	db.Callback().Update().Before("gorm:update").Register("ngo:metrics:before_update", p.metricBefore)
	db.Callback().Update().After("gorm:update").Register("ngo:metrics:after_update", p.metricAfter)

	db.Callback().Delete().Before("gorm:delete").Register("ngo:metrics:before_delete", p.metricBefore)
	db.Callback().Delete().After("gorm:delete").Register("ngo:metrics:after_delete", p.metricAfter)

	db.Callback().Row().Before("gorm:row").Register("ngo:metrics:before_row", p.metricBefore)
	db.Callback().Row().After("gorm:row").Register("ngo:metrics:after_row", p.metricAfter)

	db.Callback().Raw().Before("gorm:raw").Register("ngo:metrics:before_raw", p.metricBefore)
	db.Callback().Raw().After("gorm:raw").Register("ngo:metrics:after_raw", p.metricAfter)
}

func (p *gormMetricsPlugin) metricBefore(db *gorm.DB) {
	if !p.enable {
		return
	}
	if db == nil || db.Statement == nil || db.Statement.Context == nil {
		return
	}
	now := time.Now()
	db.InstanceSet(ngoMetricsKey, now)
}

func (p *gormMetricsPlugin) metricAfter(db *gorm.DB) {
	if !p.enable {
		return
	}
	if db == nil || db.Statement == nil || db.Statement.Context == nil {
		return
	}
	value, ok := db.InstanceGet(ngoMetricsKey)
	if !ok || value == nil {
		return
	}
	startTime, ok := value.(time.Time)
	if !ok {
		return
	}
	dsn := getDsn(db.Dialector)
	sql := sqlFilter(db.Statement.SQL.String())
	if db.Error != nil {
		fdu := float64(time.Now().Nanosecond()-startTime.Nanosecond()) / 1e6
		metricRequestError.With("dsn", dsn, "sqlString", sql).Add(1)
		metricRequestTotal.With("dsn", dsn, "sqlString", sql).Add(1)
		metricRequestDuration.With("dsn", dsn, "sqlString", sql).Set(fdu)
		metricRequestDurationRange.With("dsn", dsn, "sqlString", sql).Observe(fdu)
		metricRequestDurationAll.With("dsn", dsn, "sqlString", sql).Add(fdu)
		return
	}
	var updatedRowCount int
	var readRowCount int
	if strings.HasPrefix(strings.ToLower(sql), "select") {
		readRowCount = int(db.RowsAffected)
	} else {
		updatedRowCount = int(db.RowsAffected)
	}
	fdu := float64(time.Now().Nanosecond()-startTime.Nanosecond()) / 1e6
	metricRequestTotal.With("dsn", dsn, "sqlString", sql).Add(1)
	metricRequestDuration.With("dsn", dsn, "sqlString", sql).Set(fdu)
	metricRequestDurationRange.With("dsn", dsn, "sqlString", sql).Observe(fdu)
	metricRequestDurationAll.With("dsn", dsn, "sqlString", sql).Add(fdu)
	metricRequestEffectRow.With("dsn", dsn, "sqlString", sql).Set(float64(updatedRowCount))
	metricRequestReadRow.With("dsn", dsn, "sqlString", sql).Set(float64(readRowCount))
}
