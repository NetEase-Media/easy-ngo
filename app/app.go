package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/utils/xgo"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xlog/contrib/xzap"
	"github.com/fatih/color"

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
	status   Status
	initOnce sync.Once
}

func New() *App {
	return &App{}
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

func (app *App) Start() error {
	app.status = Starting
	ctx := context.Background()
	fs := GetFns(Starting)
	for i := range fs {
		if err := fs[i](ctx); err != nil {
			return err
		}
	}
	app.status = Running
	return nil
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
