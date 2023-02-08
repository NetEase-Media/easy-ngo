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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlFilter(t *testing.T) {
	var sql string

	sql = sqlFilter("select * from test where id in (?,?,?,?)")
	assert.Equal(t, "select * from test where id in (?)", sql)

	sql = sqlFilter("select * from test where id IN (?,?,?,?)")
	assert.Equal(t, "select * from test where id IN (?)", sql)

	sql = sqlFilter("select * from test where id IN(? , ? ,          ?,?)")
	assert.Equal(t, "select * from test where id IN (?)", sql)

	sql = sqlFilter("select * from test where id IN(? , ? ,          ?,?) and name In( ? ,? ,?)")
	assert.Equal(t, "select * from test where id IN (?) and name In (?)", sql)

	sql = sqlFilter("select * from test where id not in (?,?,?,?)")
	assert.Equal(t, "select * from test where id not in (?)", sql)

	sql = sqlFilter("select * from test where id not iN (?,?,?,?)")
	assert.Equal(t, "select * from test where id not iN (?)", sql)

	sql = sqlFilter("         insert into test in values (?,?,?,?)         ")
	assert.Equal(t, "insert into test in values (?,?,?,?)", sql)
}
