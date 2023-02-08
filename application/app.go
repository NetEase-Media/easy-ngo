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

package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/fatih/color"

	"github.com/NetEase-Media/easy-ngo/application/healthz"
	"github.com/NetEase-Media/easy-ngocation/hooks"
	"github.com/NetEase-Media/easy-ngocation/signals"
	"github.com/NetEase-Media/easy-ngocation/util"
	"github.com/NetEase-Media/easy-ngocation/util/xgo"
)

const (
	Initialize = AppStatus(iota)
	Starting
	Running
	Stopping
	Unkonwn
)

type AppStatus int

type Application struct {
	Name      string
	appStatus AppStatus
	option    *Option

	healthzServer *healthz.HealthzServer
	cycle         *util.Cycle
	smu           *sync.RWMutex
	initOnce      sync.Once
	stopOnce      sync.Once
}

func Default() *Application {
	return &Application{}
}

func (app *Application) Initialize(fns ...func() error) error {
	app.initOnce.Do(func() {

		app.printBanner()

		ctx := context.Background()
		app.cycle = util.NewCycle()
		app.smu = &sync.RWMutex{}

		app.appStatus = Initialize
		fs := hooks.GetFns(hooks.Initialize)
		for i := range fs {
			if err := fs[i](ctx); err != nil {
				util.CheckError(err)
			}
		}
	})
	return xgo.SerialUntilError(fns...)()
}

func (app *Application) Startup() error {
	app.appStatus = Starting
	ctx := context.Background()
	fns := hooks.GetFns(hooks.Start)
	if len(fns) == 0 {
		return nil
	}
	for i := range fns {
		f := fns[i]
		app.cycle.Run(func() error {
			return f(ctx)
		})
	}
	app.healthz()
	go func() {
		signals.Shutdown(func(grace bool) {
			app.Shutdown()
		})
	}()
	app.appStatus = Running
	if err := <-app.cycle.Wait(); err != nil {
		// app.option.Logger.Errorf("easy-ngo shutdown with error[%s]", err)
		return err
	}
	// app.option.Logger.Infof("shutdown easy-ngo!")
	return nil
}

func (app *Application) healthz() error {
	app.healthzServer = healthz.New()
	app.cycle.Run(app.healthzServer.Serve)
	return nil
}

func (app *Application) Shutdown() error {
	app.stopOnce.Do(func() {
		app.appStatus = Stopping
	})
	return nil
}

func (app *Application) printBanner() {

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
