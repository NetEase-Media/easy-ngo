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
	"time"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbTypeMysql = "mysql"
)

// Client
type Client struct {
	*gorm.DB
	opt     *Option
	logger  xlog.Logger
	metrics metrics.Provider
	tracer  tracer.Provider
}

func New(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*Client, error) {
	return newWithOption(opt, logger, metrics, tracer)
}

func newWithOption(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*Client, error) {
	cli := &Client{
		opt:     opt,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
	return cli, cli.initialize()
}

func (cli *Client) initialize() error {
	var cfg gorm.Config
	cfg.Logger = NewLogger(logger.Config{
		SlowThreshold: 200 * time.Millisecond,
	}, cli.logger)
	var dialector gorm.Dialector
	if cli.opt.Type == dbTypeMysql {
		dialector = mysql.Open(cli.opt.Url)
	} else {
		dialector = mysql.Open(cli.opt.Url)
	}

	db, err := gorm.Open(dialector, &cfg)
	if err != nil {
		// log.Errorf("can not be open client. msg:%s", err.Error())
		return err
	}
	myDB, err := db.DB()
	if err != nil {
		return err
	}
	myDB.SetMaxIdleConns(cli.opt.MaxIdleCons)
	myDB.SetMaxOpenConns(cli.opt.MaxOpenCons)
	myDB.SetConnMaxLifetime(cli.opt.ConnMaxLifetime)
	myDB.SetConnMaxIdleTime(cli.opt.ConnMaxIdleTime)

	if cli.metrics != nil {
		plugin := newGormMetricsPlugin(true, cli.metrics)
		db.Use(plugin)
	}

	// db.Use(newGormMetricsPlugin())
	if cli.opt.EnableTracer {
		db.Use(newGormTracerPlugin())
	}
	cli.DB = db
	return nil
}

func (cli *Client) DisConnect() error {
	if db, err := cli.DB.DB(); err == nil {
		return db.Close()
	}
	return nil
}
