package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/signals"
	"github.com/NetEase-Media/easy-ngo/utils"
	"github.com/NetEase-Media/easy-ngo/utils/xgo"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xlog/contrib/xzap"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"

	_ "github.com/NetEase-Media/easy-ngo/config/contrib/xviper"
)

const (
	Initialize Status = "Initialize"
	Starting          = "Starting"
	Running           = "Running"
	Stopping          = "Stopping"
	Online            = "Online"
	Offline           = "Offline"
	Unkonwn           = "Unkonwn"
)

type Status string

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
}

func New() *App {
	return &App{
		status:  Unkonwn,
		cycle:   utils.NewCycle(),
		smu:     &sync.RWMutex{},
		stopped: make(chan struct{}),
	}
}

func (app *App) Init(fns ...func() error) error {
	var err error
	app.initOnce.Do(func() {
		//打印logo
		app.printBanner()
		//置app状态为初始化中
		app.status = Initialize
		//初始化命令行参数
		parse()
		//初始化配置文件
		err = app.initConfig()
		if err != nil {
			return
		}
		//初始化全局日志
		err = app.initLogger()
		if err != nil {
			return
		}
		//初始化Metrics

		//初始化Tracer

		//初始化Plugins
		ctx := context.Background()
		fs := GetFns(Initialize)
		for i := range fs {
			if err := fs[i](ctx); err != nil {
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

func (app *App) initLogger() error {
	logConfig := xzap.DefaultConfig()
	if err := config.UnmarshalKey("logger", logConfig); err != nil {
		return err
	}
	xzap, _ := xzap.New(logConfig)
	xlog.WithVendor(xzap)
	return nil
}

func (app *App) initConfig() error {
	c := config.New()
	defer config.WithConfig(c)
	for _, configName := range GetConfigNames() {
		c.AddProtocol(configName)
	}
	var err error
	err = c.Init()
	if err != nil {
		return err
	}
	err = c.ReadConfig()
	if err != nil {
		return err
	}
	return nil
}

func (app *App) printBanner() {

	const banner = `
	######   ##    ####  #   #       #    #  ####   ####  
	#       #  #  #       # #        ##   # #    # #    # 
	#####  #    #  ####    #   ##### # #  # #      #    # 
	#      ######      #   #         #  # # #  ### #    # 
	#      #    # #    #   #         #   ## #    # #    # 
	###### #    #  ####    #         #    #  ####   ####  

 Welcome to easy-ngo, starting application ...
`
	fmt.Println(color.GreenString(banner))
}
