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
	"regexp"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsnMap      = make(map[string]string, 10)
	dsnMapMutex sync.RWMutex
)

var r = regexp.MustCompile(`(?i)(in)\s*\(\s*\?\s*(\s*,\s*\?)+\s*\)`)

func getConnectionCount(db *gorm.DB) int {
	sqldb, err := db.DB()
	if err == nil {
		return sqldb.Stats().OpenConnections
	}
	return 0
}

func getDsn(dialector gorm.Dialector) (dsn string) {
	var rawDsn string
	switch d := dialector.(type) {
	case nil:
	case *mysql.Dialector:
		rawDsn = d.DSN
	}
	return loadOrSaveDsn(rawDsn)
}

func loadOrSaveDsn(rawDsn string) (dsn string) {
	dsnMapMutex.RLock()
	dsn, ok := dsnMap[rawDsn]
	dsnMapMutex.RUnlock()
	if ok {
		return
	}
	dsnMapMutex.Lock()
	if dsn, ok = dsnMap[rawDsn]; !ok {
		dsn = trimDsn(rawDsn)
		dsnMap[rawDsn] = dsn
	}
	dsnMapMutex.Unlock()
	return
}

// remove user and password
// remove query
func trimDsn(dsn string) string {
	finalDsn := dsn
	start, end := 0, len(dsn)
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j := i; j >= 0; j-- {
					if dsn[j] == '@' {
						start = j + 1
						break
					}
				}
			}

			for k := i + 1; k < len(dsn); k++ {
				if dsn[k] == '?' {
					end = k
					break
				}
			}
			finalDsn = dsn[start:end]
			break
		}
	}
	return finalDsn
}

func sqlFilter(sql string) string {
	if contains(sql, " in ") || contains(sql, " in(") {
		sql = r.ReplaceAllString(sql, "${1} (?)")
	}
	return strings.TrimSpace(sql)
}

func contains(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}
