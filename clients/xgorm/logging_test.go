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
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestLogSQL(t *testing.T) {
	setupTest()
	testSqlmock.ExpectQuery("SELECT (.+) FROM `testusers`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "gender"}).FromCSVString("1,la,male"))
	var u testuser
	testClient.First(context.Background(), &u)
	assert.Equal(t, u.ID, "1")
	assert.Equal(t, u.Name, "la")
	assert.Equal(t, u.Gender, "male")
	tearDownTest()
}
