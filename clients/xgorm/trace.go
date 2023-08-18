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

// import (
// 	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/codes"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
// 	"gorm.io/gorm"
// )

// type gormTracerPlugin struct{}

// func newGormTracerPlugin() *gormTracerPlugin {
// 	return &gormTracerPlugin{}
// }

// func (p *gormTracerPlugin) Name() string {
// 	return "xgorm:tracer"
// }

// func (p *gormTracerPlugin) Initialize(db *gorm.DB) error {
// 	p.registerCallbacks(db)
// 	return nil
// }

// func (p *gormTracerPlugin) registerCallbacks(db *gorm.DB) {

// 	db.Callback().Query().Before("gorm:query").Register("ngo:tracer:before_query", p.traceBefore)
// 	db.Callback().Query().After("gorm:query").Register("ngo:tracer:after_query", p.traceAfter)

// 	db.Callback().Create().Before("gorm:create").Register("ngo:tracer:before_create", p.traceBefore)
// 	db.Callback().Create().After("gorm:create").Register("ngo:tracer:after_create", p.traceAfter)

// 	db.Callback().Update().Before("gorm:update").Register("ngo:tracer:before_update", p.traceBefore)
// 	db.Callback().Update().After("gorm:update").Register("ngo:tracer:after_update", p.traceAfter)

// 	db.Callback().Delete().Before("gorm:delete").Register("ngo:tracer:before_delete", p.traceBefore)
// 	db.Callback().Delete().After("gorm:delete").Register("ngo:tracer:after_delete", p.traceAfter)

// 	db.Callback().Row().Before("gorm:row").Register("ngo:tracer:before_row", p.traceBefore)
// 	db.Callback().Row().After("gorm:row").Register("ngo:tracer:after_row", p.traceAfter)

// 	db.Callback().Raw().Before("gorm:raw").Register("ngo:tracer:before_raw", p.traceBefore)
// 	db.Callback().Raw().After("gorm:raw").Register("ngo:tracer:after_raw", p.traceAfter)
// }

// func (p *gormTracerPlugin) traceBefore(db *gorm.DB) {
// 	// TODO enable 判断
// 	// 判断context是否存在
// 	if db == nil || db.Statement == nil || db.Statement.Context == nil {
// 		return
// 	}

// 	tr := tracer.GetTracer("gorm")
// 	db.Statement.Context, _ = tr.Start(
// 		db.Statement.Context, "gorm",
// 		tracer.WithSpanKind(tracer.SpanKindClient),
// 	)

// }

// func (p *gormTracerPlugin) traceAfter(db *gorm.DB) {
// 	// TODO enable 判断
// 	if db == nil || db.Statement == nil || db.Statement.Context == nil {
// 		return
// 	}
// 	span := tracer.SpanFromContext(db.Statement.Context)
// 	if span == nil {
// 		return
// 	}
// 	defer span.End()
// 	dsn := getDsn(db.Dialector)
// 	span.SetAttributes(
// 		semconv.DBSystemMySQL,
// 		semconv.DBConnectionStringKey.String(dsn),
// 		semconv.DBStatementKey.String(db.Statement.SQL.String()),
// 		attribute.Key("db.rowsAffected").Int64(db.Statement.RowsAffected),
// 	)
// 	span.SetStatus(codes.Ok, "success")
// 	if db.Error != nil {
// 		span.RecordError(db.Error)
// 		span.SetStatus(codes.Error, db.Error.Error())
// 	}
// }
