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

package nlog

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NetEase-Media/easy-ngo/xlog"
)

type Nlog struct {
	name  string
	flag  int // log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix| log.LUTC| log.Llongfile
	level xlog.Level
	_log  *log.Logger
}

func Default() (*Nlog, error) {
	return New(DefaultOption())
}

func New(opt *Option) (*Nlog, error) {
	alevel, err := xlog.ParseLevel(opt.Level)
	if err != nil {
		return nil, err
	}
	aflag := getFlag(opt.Flag)
	return &Nlog{name: opt.Name, flag: aflag, level: alevel, _log: log.New(os.Stdout, "", aflag)}, nil
}

func getFlag(flagStr string) int {
	flag := 0
	flagSlice := strings.Split(flagStr, "|")
	for _, oneFlagStr := range flagSlice {
		switch strings.TrimSpace(oneFlagStr) {
		case "Ldate":
			flag |= log.Ldate
		case "Ltime":
			flag |= log.Ltime
		case "Lmicroseconds":
			flag |= log.Lmicroseconds
		case "Lshortfile":
			flag |= log.Lshortfile
		case "Lmsgprefix":
			flag |= log.Lmsgprefix
		case "LUTC":
			flag |= log.LUTC
		case "Llongfile":
			flag |= log.Llongfile
		default:
			flag = 0
		}
	}
	return flag
}
func (nLog *Nlog) SetLevel(levelStr string) error {
	alevel, err := xlog.ParseLevel(levelStr)
	if err != nil {
		return err
	}
	nLog.level = alevel
	return nil
}
func (nLog *Nlog) GetLevel() string {
	return nLog.level.String()
}
func (nLog *Nlog) SetFlags(flagStr string) {
	nLog.flag = getFlag(flagStr)
	nLog._log.SetFlags(nLog.flag)
}
func (nLog *Nlog) GetFlags() int {
	return nLog.flag
}
func (nLog *Nlog) GetName() string {
	return nLog.name
}
func (nLog *Nlog) String() string {
	return "&{" + nLog.name + " " + nLog.level.String() + " " + strconv.FormatInt(int64(nLog.flag), 2) + "}"
}
func (nLog *Nlog) Debugf(format string, fields ...interface{}) {
	if nLog.level <= xlog.DebugLevel {
		nLog._log.Printf("[\x1b[0;35m"+xlog.Level.CapitalString(xlog.DebugLevel)+"\x1b[0m] "+format, fields...)
	}
}

func (nLog *Nlog) Infof(format string, fields ...interface{}) {
	if nLog.level <= xlog.InfoLevel {
		nLog._log.Printf("[\x1b[0;34m"+xlog.Level.CapitalString(xlog.InfoLevel)+"\x1b[0m] "+format, fields...)
	}
}

func (nLog *Nlog) Warnf(format string, fields ...interface{}) {
	if nLog.level <= xlog.WarnLevel {
		nLog._log.Printf("[\x1b[0;33m"+xlog.Level.CapitalString(xlog.WarnLevel)+"\x1b[0m] "+format, fields...)
	}
}

func (nLog *Nlog) Errorf(format string, fields ...interface{}) {
	if nLog.level <= xlog.ErrorLevel {
		nLog._log.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.ErrorLevel)+"\x1b[0m] "+format, fields...)
	}
}

func (nLog *Nlog) DPanicf(format string, fields ...interface{}) {
	if nLog.level <= xlog.DPanicLevel {
		nLog._log.Printf("["+xlog.Level.CapitalString(xlog.InfoLevel)+"] "+format, fields...)
	}
}

func (nLog *Nlog) Panicf(format string, fields ...interface{}) {
	if nLog.level <= xlog.PanicLevel {
		nLog._log.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.PanicLevel)+"\x1b[0m] "+format, fields...)
	}
}

func (nLog *Nlog) Fatalf(format string, fields ...interface{}) {
	if nLog.level <= xlog.FatalLevel {
		nLog._log.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.FatalLevel)+"\x1b[0m] "+format, fields...)
	}
}

func (nLog *Nlog) Level() xlog.Level {
	return nLog.level
}
