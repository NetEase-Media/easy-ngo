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

// import (
// 	"context"
// 	"testing"

// 	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/schema"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/agiledragon/gomonkey/v2"
// 	"github.com/stretchr/testify/assert"
// )

// type user struct {
// 	Id   int64
// 	Name string
// }

// func (user) TableName() string {
// 	return "test"
// }

// func TestClient(t *testing.T) {
// 	patches := gomonkey.ApplyFunc(mysql.Open, func(dsn string) gorm.Dialector {
// 		db, _, _ := sqlmock.New()
// 		return mysql.New(mysql.Config{
// 			DSN:                       dsn,
// 			SkipInitializeWithVersion: true,
// 			Conn:                      db,
// 		})
// 	})
// 	defer patches.Reset()

// 	c, err := newWithOption(&Option{
// 		Name:            "test",
// 		Type:            "mysql",
// 		Url:             "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
// 		MaxIdleCons:     10,
// 		MaxOpenCons:     10,
// 		ConnMaxLifetime: 1000,
// 		ConnMaxIdleTime: 10,
// 	}, &xfmt.XFmt{}, nil, nil)
// 	assert.Nil(t, err)
// 	c.NamingStrategy = schema.NamingStrategy{
// 		SingularTable: true,
// 	}
// 	var tb = user{}
// 	ctx := context.Background()
// 	err = c.Create(&tb).Error
// 	assert.Nil(t, err)
// 	err = c.Find(ctx, &tb).Error
// 	assert.Nil(t, err)
// }
