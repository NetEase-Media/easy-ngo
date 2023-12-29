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

package app

import (
	"context"
	"sync"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xmetrics"
	"github.com/NetEase-Media/easy-ngo/xtracer"

	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/signals"
	"github.com/NetEase-Media/easy-ngo/utils"
	"github.com/NetEase-Media/easy-ngo/utils/xgo"
	"github.com/NetEase-Media/easy-ngo/xlog/contrib/xstdout"
	"github.com/NetEase-Media/easy-ngo/xlog/contrib/xzap"
	"github.com/NetEase-Media/easy-ngo/xmetrics/contrib/xprometheus"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
)

const (
	Initialize Status = "Initialize"
	Starting   Status = "Starting"
	Running    Status = "Running"
	Stopping   Status = "Stopping"
	Online     Status = "Online"
	Offline    Status = "Offline"
	Unkonwn    Status = "Unkonwn"
)

type Status string

type BaseConfigKey string

const (
	LoggerConfigKey  BaseConfigKey = "logger"
	TracerConfigKey  BaseConfigKey = "tracer"
	MetricsConfigKey BaseConfigKey = "metrics"
)

type App struct {
	status Status
	//保证初始化只执行一次
	initOnce sync.Once
	//保证启动只执行一次
	startOnce sync.Once
	//保证停止只执行一次
	stopOnce sync.Once

	cycle   *utils.Cycle
	smu     *sync.RWMutex
	stopped chan struct{}

	// if set true, we will init tracer plugin
	enableTracer bool
	// if set true, we will init metric plugin
	enableMetrics bool
}

func New() *App {
	return &App{
		status:  Unkonwn,
		cycle:   utils.NewCycle(),
		smu:     &sync.RWMutex{},
		stopped: make(chan struct{}),
	}
}

func (app *App) EnableTracer() *App {
	app.enableTracer = true
	return app
}

func (app *App) EnableMetrics() *App {
	app.enableMetrics = true
	return app
}

func (app *App) Init(fns ...func() error) error {
	var err error
	app.initOnce.Do(func() {
		//init standard output and print logo
		app.initStdoutAndPrintBanner()
		//set app status with Initialize
		app.status = Initialize
		//初始化命令行参数
		parse()
		//初始化配置文件
		err = app.initConfig()
		if err != nil {
			xlog.Errorf("init config error: %v", err.Error())
			return
		}
		//初始化全局日志
		err = app.initLogger()
		if err != nil {
			xlog.Errorf("init logger error: %v", err.Error())
			return
		}
		//初始化Metrics
		err = app.initMetrics()
		if err != nil {
			xlog.Errorf("init metrics error: %v", err.Error())
			return
		}
		//初始化Tracer
		err = app.initTracer()
		if err != nil {
			xlog.Errorf("init tracer error: %v", err.Error())
			return
		}
		//初始化Plugins
		ctx := context.Background()
		fs := GetFns(Initialize)
		for i := range fs {
			if err := fs[i](ctx); err != nil {
				xlog.Errorf("init plugins error: %v", err.Error())
				return
			}
		}
		err = xgo.SerialUntilError(fns...)()
	})
	return err
}

func (app *App) Start(fns ...func() error) error {
	var err error
	app.startOnce.Do(func() {
		//如果App状态为Unkonwn，说明没有执行过Init，需要先执行Init
		if app.status == Unkonwn {
			if err = app.Init(fns...); err != nil {
				return
			}
		}
		app.status = Starting
		app.cycle.Run(app.startPlugins)
		app.waitSignals()
		app.status = Running
		xlog.Infof("easy-ngo start success!")
		if err := <-app.cycle.Wait(); err != nil {
			xlog.Errorf("easy-ngo shutdown with error[%s]", err.Error())
			return
		}
		xlog.Infof("shutdown easy-ngo!")
	})
	return err
}

func (app *App) startPlugins() error {
	app.smu.Lock()
	defer app.smu.Unlock()
	var eg errgroup.Group
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		<-app.stopped
		cancel()
	}()
	fs := GetFns(Starting)
	for _, f := range fs {
		eg.Go(func() (err error) {
			err = f(ctx)
			return
		})
	}
	return eg.Wait()
}

func (app *App) waitSignals() {
	app.smu.Lock()
	defer app.smu.Unlock()
	signals.Shutdown(func(grace bool) {
		xlog.Infof("easy-ngo Stopping!")
		_ = app.Shutdown()
		xlog.Infof("easy-ngo Stopped!")
	})
}

func (app *App) Shutdown() (err error) {
	app.stopOnce.Do(func() {
		var eg errgroup.Group
		var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
		app.stopped <- struct{}{}
		fs := GetFns(Stopping)
		for _, f := range fs {
			eg.Go(func() (err error) {
				err = f(ctx)
				return
			})
		}
		app.cycle.Close()
	})
	return
}

func (app *App) initTracer() error {
	if !app.enableTracer {
		return nil
	}
	tracerConfig := xtracer.DefaultConfig()
	if err := config.UnmarshalKey(string(TracerConfigKey), tracerConfig); err != nil {
		return err
	}
	provider := xtracer.New(tracerConfig)
	xtracer.WithVendor(provider)
	return nil
}

func (app *App) initMetrics() error {
	if !app.enableMetrics {
		return nil
	}
	metricsConfig := xprometheus.DefaultConfig()
	if err := config.UnmarshalKey(string(MetricsConfigKey), metricsConfig); err != nil {
		return err
	}
	provider := xprometheus.NewProvider(metricsConfig)
	xmetrics.WithVendor(provider)
	server := xprometheus.NewServer(metricsConfig)
	go server.Start()
	return nil
}

func (app *App) initLogger() (err error) {
	if !config.Exists(string(LoggerConfigKey)) {
		return nil
	}
	var logConfig *xzap.Config = xzap.DefaultConfig()
	if err = config.UnmarshalKey(string(LoggerConfigKey), logConfig); err != nil {
		return err
	}
	if logConfig == nil {
		return nil
	}
	logger, err := xzap.New(logConfig)
	if err != nil {
		return err
	}
	xlog.WithVendor(logger)
	return nil
}

func (app *App) initConfig() error {
	conf := config.New()
	if err := conf.Init(GetConfigProtocols()...); err != nil {
		return err
	}
	config.WithConfig(conf)
	return nil
}

func (app *App) initStdoutAndPrintBanner() {
	var logger xlog.Logger = xstdout.New()
	xlog.WithVendor(logger)

	const banner = `
	######   ##    ####  #   #       #    #  ####   ####  
	#       #  #  #       # #        ##   # #    # #    # 
	#####  #    #  ####    #   ##### # #  # #      #    # 
	#      ######      #   #         #  # # #  ### #    # 
	#      #    # #    #   #         #   ## #    # #    # 
	###### #    #  ####    #         #    #  ####   ####  

 Welcome to easy-ngo, starting application ...
`
	xlog.Infof(color.GreenString(banner))
}
