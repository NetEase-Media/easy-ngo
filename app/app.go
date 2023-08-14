package app

import (
	"fmt"
	"sync"

	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/fatih/color"
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
	app.initOnce.Do(func() {
		//打印logo
		app.printBanner()
		//置app状态为初始化中
		app.status = Initialize
		//初始化命令行参数
		parse()
		//初始化配置文件
		var err error
		err = app.initConfig()
		if err != nil {
			panic(err)
		}
		//初始化全局日志

		//初始化Metrics

		//初始化Tracer

		//初始化Plugins
	})
	return nil
}

func (app *App) initConfig() error {
	c := config.New()
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
