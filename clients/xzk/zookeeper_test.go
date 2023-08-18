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

package xzk

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
)

const (
	// ADDR = "127.0.0.1:2181"
	ADDR = "kafka"
	NAME = "zktest"
)

func TestCreateNode_00(t *testing.T) {
	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}
	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")
	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}
}

func TestOptionsProxy_GetConnState_flag0(t *testing.T) {
	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}
	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	time.Sleep(5 * time.Second)
	state := c.GetConnState()
	assert.Equal(t, "StateHasSession", state, "获取连接状态失败")
}

func TestOptionsProxy_GetSessionId_flag0(t *testing.T) {
	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}
	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	s := c.GetSessionId()
	log.Println(s)
	assert.Equal(t, s, c.GetSessionId(), "获取sessionID失败")
}

func TestOptionsProxy_CreateNode_flag0(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点
	dd, _ := c.Exist("/cc")
	if dd == true {
		ddd := c.Delete("/cc")
		assert.Equal(t, nil, ddd, "删除节点失败")
	}
	path, b := c.CreateNode("/cc", PERSISTENT, "aa")
	assert.Equal(t, nil, b, "创建失败")
	assert.Equal(t, "/cc", path, "创建失败")
	if b == nil {
		c.Delete("/cc")
	}
}

func TestOptionsProxy_CreateNodeWithAcls_flag1(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr: []string{ADDR},
		// SessionTimeout: time.Second * 5,             // SessionTimeout的默认值为 5*time.Second
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")

	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	acls := zk.WorldACL(zk.PermRead)
	path := "/pushTest/test5"
	realPath, bb := c.CreateNodeWithAcls(path, EPHEMERAL_SEQUENTIAL, acls, "eee")
	assert.Equal(t, nil, bb, "创建节点失败")

	da, ee := c.GetData(realPath)
	assert.Equal(t, nil, ee, "获取节点值失败")
	assert.Equal(t, "eee", da, "获取节点值失败")
	eee := c.SetData(realPath, "eiei")
	assert.NotNil(t, eee)
}
func TestOptionsProxy_CreateNode_flag1(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr: []string{ADDR},
		// SessionTimeout: time.Second * 5,             //  SessionTimeout的默认值为 5*time.Second
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点
	dd, _ := c.Exist("/cc")
	if dd == true {
		ddd := c.Delete("/cc")
		assert.Equal(t, nil, ddd, "删除节点失败")
	}
	_, b := c.CreateNode("/cc", EPHEMERAL_SEQUENTIAL, "aa")
	assert.Equal(t, nil, b, "创建失败")
	/*dd,_ = c.Exist("/cc")
	assert.Equal(t, false, dd, "节点未自动失效")*/
}

func TestOptionsProxy_CreateNode_flag2(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点(带序号)
	_, bb := c.CreateNode("/testzk_ss1s_", PERSISTENT_SEQUENTIAL, "flag2")
	assert.Equal(t, nil, bb, "创建失败")
	str, _, _ := c.Conn.Children("/")
	for i := range str {
		if strings.Contains(str[i], "testzk_ss1s_") {
			c.Delete("/" + str[i])
		}
	}
}

func TestOptionsProxy_CreateNodeParent_flag0(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	dd, _ := c.Exist("/parentNode1/cc/aa")
	if dd == true {
		b := c.Delete("/parentNode1/cc/aa")
		bb := c.Delete("/parentNode1/cc")
		bbb := c.Delete("/parentNode1")
		assert.Equal(t, true, bb, "创建失败")
		assert.Equal(t, true, bbb, "创建失败")
		assert.Equal(t, true, b, "创建失败")
	}
	fmt.Println(111)
	// 创建节点
	path, bb := c.CreateNodeParent("/parentNode1/cc/aa", PERSISTENT, "flag4")
	assert.Equal(t, nil, bb, "创建失败")
	assert.Equal(t, "/parentNode1/cc/aa", path, "创建失败")
	s, _ := c.Exist("/parentNode1")
	ss, _ := c.Exist("/parentNode1/cc")
	sss, _ := c.Exist("/parentNode1/cc/aa")
	assert.Equal(t, true, s, "创建父节点失败")
	assert.Equal(t, true, ss, "创建节点失败")
	assert.Equal(t, true, sss, "创建节点失败")

	_, bb1 := c.CreateNodeParent("/parentNode1/cc/aa/bb", PERSISTENT, "flag4")
	assert.Equal(t, nil, bb1, "创建失败")

	ressss := c.Delete("/parentNode1/cc/aa/bb")
	resss := c.Delete("/parentNode1/cc/aa")
	res := c.Delete("/parentNode1/cc")
	ress := c.Delete("/parentNode1")
	assert.Nil(t, res)
	assert.Nil(t, ress)
	assert.Nil(t, resss)
	assert.Nil(t, ressss)
}

func TestOptionsProxy_CreateNodeParent_flag1(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	dd, _ := c.Exist("/parentNode2")
	if dd == true {
		bbb := c.Delete("/parentNode2")
		assert.Equal(t, true, bbb, "创建失败")
	}
	path, bb := c.CreateNodeParent("/parentNode202214", PERSISTENT, "flag4")
	assert.Equal(t, nil, bb, "创建失败")
	assert.Equal(t, "/parentNode202214", path, "创建失败")
	ress := c.Delete("/parentNode202214")
	assert.Nil(t, ress)
}

func TestOptionsProxy_CreateNodeParent_flag2(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点(带序号)
	path, bb := c.CreateNodeParent("/parentNode3/cc", PERSISTENT_SEQUENTIAL, "flag4")
	assert.Equal(t, nil, bb, "创建失败")

	_, bb1 := c.CreateNodeParent(path+"/aa", PERSISTENT_SEQUENTIAL, "flag4")
	assert.Equal(t, nil, bb1, "创建失败")

	pp := strings.Split(path, "/")
	ee, _ := c.Exist("/" + pp[1])
	assert.Equal(t, true, ee, "创建失败")

}

func TestOptionsProxy_CreateNodeParent_flag3(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点(带序号)
	path, bb := c.CreateNodeParent("/parentNode4/cc", EPHEMERAL_SEQUENTIAL, "flag4")
	assert.Equal(t, nil, bb, "创建失败")
	pp := strings.Split(path, "/")
	ee, _ := c.Exist("/" + pp[1])
	assert.Equal(t, true, ee, "创建失败")
}

func TestOptionsProxy_Exist(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 5,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	path, er := c.CreateNode("/existTest", PERSISTENT_SEQUENTIAL, "aa")
	path1, _ := c.CreateNode(path+"/child", PERSISTENT_SEQUENTIAL, "bb")
	assert.Equal(t, nil, er, "创建失败")
	isExist, _ := c.Exist(path)
	assert.Equal(t, true, isExist, "节点判断失败")
	isExist1, _ := c.Exist(path1)
	assert.Equal(t, true, isExist1, "节点判断失败")
}

func TestOptionsProxy_SetData_pathErr(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点(带序号)
	path, bb := c.CreateNode("/testzk_ss1s_", PERSISTENT_SEQUENTIAL, "flag2")
	log.Println(path)
	assert.Equal(t, nil, bb, "创建失败")
	q := c.SetData(path, "testzk_ss1s_setData")
	assert.Equal(t, nil, q, "设置节点值失败")
}

func TestOptionsProxy_SetData(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点(带序号)
	_, bb := c.CreateNode("/testzk_ss1s_", PERSISTENT_SEQUENTIAL, "flag2")
	assert.Equal(t, nil, bb, "创建失败")
	str, _, _ := c.Conn.Children("/")
	for i := range str {
		if strings.Contains(str[i], "testzk_ss1s_") {
			q := c.SetData("/"+str[i], "testzk_ss1s_setData")
			assert.Equal(t, nil, q, "设置节点值失败")
			c.Delete("/" + str[i])
			break
		}
	}
}

func TestOptionsProxy_GetNodeData(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	dd, _ := c.Exist("/testzk_ss0s")
	if dd == false {
		_, bb := c.CreateNode("/testzk_ss0s", PERSISTENT, "flag0")
		assert.Equal(t, nil, bb, "创建失败")
	}
	q := c.SetData("/testzk_ss0s", "testzk_ss1s_setData")
	assert.Equal(t, nil, q, "设置节点值失败")
	data, err := c.GetData("/testzk_ss0s")
	assert.Equal(t, nil, err, "设置节点值失败")
	assert.Equal(t, "testzk_ss1s_setData", data, "设置节点值失败")
}

func TestOptionsProxy_GetChildren(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	rootPath, _ := c.CreateNode("/getChildrenTest", PERSISTENT_SEQUENTIAL, "aa")
	childPath1, errr := c.CreateNode(rootPath+"/a", PERSISTENT, "a")
	assert.Equal(t, nil, errr, "创建失败")
	childPath2, _ := c.CreateNode(rootPath+"/b", PERSISTENT, "b")
	childSlice, err := c.GetChildren(rootPath)
	assert.Equal(t, nil, err, "获取子节点失败")
	assert.Equal(t, childPath1, childSlice[0], "获取子节点失败")
	assert.Equal(t, childPath2, childSlice[1], "获取子节点失败")
}

func TestOptionsProxy_Delete(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	dd, _ := c.Exist("/testzk_ss0s")
	if dd == true {
		c.Delete("/testzk_ss0s")
	}
	// 创建节点
	_, bb := c.CreateNode("/testzk_ss0s", PERSISTENT, "flag0")
	assert.Equal(t, nil, bb, "创建失败")
	res := c.Delete("/testzk_ss0s")
	assert.Nil(t, res)
}

func TestOptionsProxy_Delete_PathErr(t *testing.T) {

	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点
	dd, _ := c.Exist("/testzk_ss0s")
	if dd == true {
		c.Delete("/testzk_ss0s")
	}
	_, bb := c.CreateNode("/testzk_ss0s", PERSISTENT, "flag0")
	assert.Equal(t, nil, bb, "创建失败")
	res := c.Delete("/testzk_ss0s")
	assert.Equal(t, nil, res, "删除失败")
}

func TestAll(t *testing.T) {

	o := Config{
		Name: "mm",
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}

	c, err := New(&o)
	// err := Init(opts)
	assert.Equal(t, nil, err, "初始化失败")

	// 创建节点
	dd, _ := c.Exist("/cc")
	if dd == true {
		ddd := c.Delete("/cc")
		assert.Equal(t, true, ddd, "创建节点失败")
	}
	_, b := c.CreateNode("/cc", PERSISTENT, "aa")
	assert.Equal(t, nil, b, "创建失败")
	if b == nil {
		c.Delete("/cc")
	}
	_, b = c.CreateNode("/cc", EPHEMERAL, "aa1")
	assert.Equal(t, nil, b, "创建失败")

	// 创建节点(带序号)
	_, bb := c.CreateNode("/ss1s_", PERSISTENT_SEQUENTIAL, "aa")
	assert.Equal(t, nil, bb, "创建失败")

	// 检查节点是否存在
	d, _ := c.Exist("/cc")
	assert.Equal(t, true, d, "该节点不存在")

	// 设置节点值
	dddd, _ := c.Exist("/bb")
	if dddd == false {
		_, sign := c.CreateNode("/bb", PERSISTENT, "ss")
		assert.Equal(t, nil, sign, "")
	}
	q := c.SetData("/bb", "ahs")
	assert.Equal(t, nil, q, "设置节点值失败")

	// 获取节点值
	e, _ := c.Exist("/ee")
	if e == false {
		_, sign := c.CreateNode("/ee", PERSISTENT, "ss")
		assert.Equal(t, nil, sign, "")
	}
	eee := c.SetData("/ee", "ss")
	assert.Equal(t, nil, eee, "设置节点值失败")
	eeee, _ := c.GetData("/ee")
	assert.Equal(t, "ss", eeee, "获取节点值失败")

	// 删除节点
	ds, _ := c.Exist("/qq")
	if ds == false {
		_, sign1 := c.CreateNode("/qq", PERSISTENT, "ss")
		assert.Equal(t, nil, sign1, "")
	}
	p := c.Delete("/qq")
	assert.Equal(t, nil, p, "删除节点失败")

}

// 监视子节点
func TestZkChildNodeWatcher(t *testing.T) {
	o := Config{
		Name: "cc",
		Addr: []string{ADDR},
		//Addr:           []string{"localhost:2181"},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}
	path := "/pushTest/testWatchChildNode1"

	realPath, bb := c.CreateNode(path, PERSISTENT_SEQUENTIAL, "aa")
	assert.Equal(t, nil, bb, "创建失败")
	c.WatchChildren(realPath, func(respChan <-chan *WatchChildrenResponse) {
		for r := range respChan {
			for k, _ := range r.ChildrenChangeInfo {
				fmt.Println(r.ChildrenChangeInfo[k].Path+"  "+r.ChildrenChangeInfo[k].OperateType.String()+" err:", r.Err)
			}
		}
	})
	time.Sleep(1 * time.Second)

	pathc := realPath + "/bbbb"
	_, b := c.CreateNode(pathc, PERSISTENT, "aa")

	pathcc := realPath + "/cccb"
	_, d := c.CreateNode(pathcc, PERSISTENT, "aa")

	bbb := c.Delete(pathc)

	pathccc := realPath + "/ddd"
	_, dd := c.CreateNode(pathccc, PERSISTENT, "aa")

	pathcccc := realPath + "/dddddddd"
	_, ddd := c.CreateNode(pathcccc, PERSISTENT, "aa")

	pathccccc := realPath + "/ddddd"
	_, dddd := c.CreateNode(pathccccc, PERSISTENT, "aa")

	ddddbbbb := c.Delete(pathcc)
	fmt.Println("dddddd")
	_, dddd1 := c.CreateNode(pathcc, PERSISTENT, "aa")
	assert.Equal(t, nil, dddd1, "删除子节点失败")
	// time.Sleep(120 * time.Second) 				// 关闭打开VPN   测试断开连接 及相应的检测机制
	pathcccccc := realPath + "/dddddc"
	_, ddddd := c.CreateNode(pathcccccc, PERSISTENT, "aa")

	pathccccccc := realPath + "/dddddcc"
	_, dddddd := c.CreateNode(pathccccccc, PERSISTENT, "aa")

	pathcccccccc := realPath + "/dddddccc"
	_, ddddddd := c.CreateNode(pathcccccccc, PERSISTENT, "aa")
	time.Sleep(1 * time.Second)

	assert.Equal(t, nil, b, "创建子节点失败")
	assert.Equal(t, nil, d, "创建子节点失败")
	assert.Equal(t, nil, bbb, "删除子节点失败")
	assert.Equal(t, nil, dd, "创建子节点失败")
	assert.Equal(t, nil, ddd, "创建子节点失败")
	assert.Equal(t, nil, dddd, "创建子节点失败")
	assert.Equal(t, nil, ddddd, "创建子节点失败")
	assert.Equal(t, nil, dddddd, "创建子节点失败")
	assert.Equal(t, nil, ddddddd, "创建子节点失败")
	assert.Equal(t, nil, ddddbbbb, "删除子节点失败")
	childAll, _ := c.GetChildren(realPath)
	fmt.Println("所有子节点: ", childAll)
	for _, v := range childAll {
		d111 := c.Delete(v)
		assert.Equal(t, nil, d111, "删除失败 path:"+v)
		fmt.Println("删除节点: " + v)
	}
	childA, _ := c.GetChildren(realPath)

	assert.Equal(t, 0, len(childA), realPath+" 还有子节点")
	if len(childA) == 0 {
		dg := c.Delete(realPath)
		assert.Equal(t, nil, dg, "删除失败 path:"+realPath)
	}
}

func TestZkChildNodeWatcher_02(t *testing.T) {

	o := Config{
		Name: "cc11",
		Addr: []string{ADDR},
		// Addr:           []string{"localhost:2181"},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	path := "/pushTest/testWatchChildNode2"

	realPath, bb := c.CreateNode(path, PERSISTENT_SEQUENTIAL, "aa")
	assert.Equal(t, nil, bb, "创建失败")
	c.WatchChildren(realPath, func(respChan <-chan *WatchChildrenResponse) {
		for r := range respChan {
			for k, _ := range r.ChildrenChangeInfo {
				fmt.Println(r.ChildrenChangeInfo[k].Path+"  "+r.ChildrenChangeInfo[k].OperateType.String()+" err:", r.Err)
			}
		}
	})
	pathc := realPath + "/bbbb"
	_, b := c.CreateNode(pathc, PERSISTENT, "aa")
	assert.Equal(t, nil, b, "创建子节点失败")
}

func TestZkChildNodeWatcher_03(t *testing.T) {

	o := Config{
		Name: "cc",
		Addr: []string{ADDR},
		//Addr:           []string{"localhost:2181"},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	path := "/pushTest/testWatchChildNode3"

	realPath, bb := c.CreateNode(path, EPHEMERAL_SEQUENTIAL, "aa")
	assert.Equal(t, nil, bb, "创建失败")
	c.WatchChildren(realPath, func(respChan <-chan *WatchChildrenResponse) {
		for r := range respChan {
			for k, _ := range r.ChildrenChangeInfo {
				fmt.Println(r.ChildrenChangeInfo[k].Path+"  "+r.ChildrenChangeInfo[k].OperateType.String()+" err:", r.Err)
			}
		}
	})
	pathc := realPath + "/bbbb"
	_, b := c.CreateNode(pathc, PERSISTENT, "aa")
	assert.NotNil(t, b)
}

// 测试临时节点的监听： 断开重连、断开恢复临时节点、断开恢复对临时节点的监听
func TestZkNodeWatcher_EphemeralNode_EPHEMERAL_04(t *testing.T) {

	o := Config{
		Name:           "cc1",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	path := "/pushTest/test2022_1_25"

	realPath, bb := c.CreateNode(path, EPHEMERAL, "")
	assert.Equal(t, nil, bb, "创建失败")
	backchannel, _ := c.StartLi()
	go func(backchannel <-chan Event) {
		for kk := range backchannel {
			fmt.Println(kk)
		}
	}(backchannel)
	sign := c.SetData(realPath, "bb")
	assert.Equal(t, nil, sign, "修改失败")
	// c.WatchNode(realPath, func(respChan <-chan *WatchNodeResponse) {
	// 	for d := range respChan {
	// 		fmt.Println("正常监听输出 "+path+"newData11 "+d.NodeChangeInfo.OperateType.String()+"  "+d.NodeChangeInfo.OldData+"  "+d.NodeChangeInfo.NewData, "  ", d.Err)
	// 	}
	// })

	isE, err1 := c.Exist(realPath)
	assert.Equal(t, nil, err1, "判断是否存在失败")
	assert.Equal(t, true, isE, "不存在")
	sign1 := c.SetData(realPath, "CC")
	assert.Equal(t, nil, sign1, "修改失败")
	// time.Sleep(120 * time.Second) // 断开又打开VPN

	isEe, err1e := c.Exist(realPath)
	assert.Equal(t, nil, err1e, "判断是否存在失败")
	assert.Equal(t, true, isEe, "不存在")
	daa, _ := c.GetData(realPath)
	assert.Equal(t, "CC", daa, "断开重连创建节点值失败")

	sign11 := c.SetData(realPath, "DD")
	assert.Equal(t, nil, sign11, "修改失败")
	sign111 := c.SetData(realPath, "EE")
	assert.Equal(t, nil, sign111, "修改失败")
	for i := 0; i < 100; i++ {
		sign := c.SetData(realPath, "realPath_"+strconv.Itoa(i))
		assert.Equal(t, nil, sign, "修改失败")
	}
	sign1111 := c.SetData(realPath, "FF")
	assert.Equal(t, nil, sign1111, "修改失败")
	daaa, _ := c.GetData(realPath)
	assert.Equal(t, "FF", daaa, "断开重连创建节点值失败")
	fmt.Println("最终后data: ", daaa)
	time.Sleep(1 * time.Second)

	path1 := "/pushTest/test2022_1_252"

	path11, bbC := c.CreateNode(path1, EPHEMERAL, "")
	assert.Equal(t, nil, bbC, "创建失败")

	sign11112 := c.SetData(path11, "FF")
	assert.Equal(t, nil, sign11112, "修改失败")
}

// 对临时节点监听实时性测试
func TestZkNodeWatcher_EphemeralNode_EPHEMERAL_05(t *testing.T) {

	o := Config{
		Name:           "cc3",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	path := "/pushTest/test_01"

	realPath, bb := c.CreateNode(path, EPHEMERAL_SEQUENTIAL, "")
	assert.Equal(t, nil, bb, "创建失败")
	backchannel, _ := c.StartLi()
	go func(backchannel <-chan Event) {
		for kk := range backchannel {
			fmt.Println(kk)
		}
	}(backchannel)
	sign := c.SetData(realPath, "bb")
	assert.Equal(t, nil, sign, "修改失败")
	c.WatchNode(realPath, func(respChan <-chan *WatchNodeResponse) {
		for d := range respChan {
			fmt.Println("正常监听输出 "+path+"newData11 "+d.NodeChangeInfo.OperateType.String()+"  "+d.NodeChangeInfo.OldData+"  "+d.NodeChangeInfo.NewData, "  ", d.Err)
		}
	})
	sign11 := c.SetData(realPath, "DD")
	assert.Equal(t, nil, sign11, "修改失败")
	// time.Sleep(120 * time.Second) // 关闭打开VPN
	fmt.Println(9999)
	sign111 := c.SetData(realPath, "EE")
	assert.Equal(t, nil, sign111, "修改失败")
	time.Sleep(1 * time.Second)
}

func TestZkNodeWatcher_Node_PERSISTENT_05(t *testing.T) {

	o := Config{
		Name:           "cc4",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	path := "/pushTest/test_02"

	backchannel, _ := c.StartLi()
	go func(backchannel <-chan Event) {
		for kk := range backchannel {
			fmt.Println(kk)
		}
	}(backchannel)

	realPath, bb := c.CreateNode(path, PERSISTENT_SEQUENTIAL, "")
	assert.Equal(t, nil, bb, "创建失败")

	sign := c.SetData(realPath, "bb")
	assert.Equal(t, nil, sign, "修改失败")
	c.WatchNode(realPath, func(respChan <-chan *WatchNodeResponse) {
		for d := range respChan {
			fmt.Println("11正常监听输出 "+path+"newData11 "+d.NodeChangeInfo.OperateType.String()+"  "+d.NodeChangeInfo.OldData+"  "+d.NodeChangeInfo.NewData, "  ", d.Err)
		}
	})

	c.WatchNode(realPath, func(respChan <-chan *WatchNodeResponse) {
		for d := range respChan {
			fmt.Println("22正常监听输出 "+path+"newData11 "+d.NodeChangeInfo.OperateType.String()+"  "+d.NodeChangeInfo.OldData+"  "+d.NodeChangeInfo.NewData, "  ", d.Err)
		}
	})

	isE, err1 := c.Exist(realPath)
	assert.Equal(t, nil, err1, "判断是否存在失败")
	assert.Equal(t, true, isE, "不存在")
	fmt.Println(9991111)
	sign1 := c.SetData(realPath, "CC")
	assert.Equal(t, nil, sign1, "修改失败")
	// time.Sleep(120 * time.Second)            // 关闭 打开VPN
	fmt.Println(88888)
	sign11 := c.SetData(realPath, "DD")
	assert.Equal(t, nil, sign11, "修改失败")
	sign111 := c.SetData(realPath, "EE")
	assert.Equal(t, nil, sign111, "修改失败")
	SSSS := c.Delete(realPath)
	assert.Equal(t, nil, SSSS, "删除失败")

	realPath, bb = c.CreateNode(realPath, PERSISTENT, "")
	assert.Equal(t, nil, bb, "创建失败")
	time.Sleep(2 * time.Second)
	fmt.Println(77777)
	sign1111 := c.SetData(realPath, "FF")
	assert.Equal(t, nil, sign1111, "修改失败")
	time.Sleep(1 * time.Second)
}

func TestZkNodeWatcher_Node_PERSISTENT_06(t *testing.T) {

	o := Config{
		Name:           "cc5",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

	path := "/pushTest/test_03"

	backchannel, _ := c.StartLi()
	go func(backchannel <-chan Event) {
		for kk := range backchannel {
			fmt.Println(kk)
		}
	}(backchannel)

	realPath, bb := c.CreateNode(path, PERSISTENT_SEQUENTIAL, "")
	assert.Equal(t, nil, bb, "创建失败")

	sign := c.SetData(realPath, "bb")
	assert.Equal(t, nil, sign, "修改失败")
	c.WatchNode(realPath, func(respChan <-chan *WatchNodeResponse) {
		for d := range respChan {
			fmt.Println("11正常监听输出 "+path+"newData11 "+d.NodeChangeInfo.OperateType.String()+"  "+d.NodeChangeInfo.OldData+"  "+d.NodeChangeInfo.NewData, "  ", d.Err)
		}
	})

	isE, err1 := c.Exist(realPath)
	assert.Equal(t, nil, err1, "判断是否存在失败")
	assert.Equal(t, true, isE, "不存在")
	fmt.Println(9991111)
	sign1 := c.SetData(realPath, "CC")
	assert.Equal(t, nil, sign1, "修改失败")
	// time.Sleep(120 * time.Second) // 关闭 打开VPN
	ddd, e11 := c.GetData(realPath)
	assert.Equal(t, nil, e11, "获取节点值失败")
	fmt.Println("断开重连后的data: ", ddd)
	sign11 := c.SetData(realPath, "DD")
	assert.Equal(t, nil, sign11, "修改失败")
	sign111 := c.SetData(realPath, "EE")
	assert.Equal(t, nil, sign111, "修改失败")
	SSSS := c.Delete(realPath)
	assert.Equal(t, nil, SSSS, "删除失败")
	time.Sleep(5 * time.Second)
	fmt.Println(999)
	realPath, bb = c.CreateNode(realPath, PERSISTENT, "")
	assert.Equal(t, nil, bb, "创建失败")
	time.Sleep(2 * time.Second)
	fmt.Println(77777)
	sign1111 := c.SetData(realPath, "FF")
	assert.Equal(t, nil, sign1111, "修改失败")
	time.Sleep(1 * time.Second)
	fmt.Println(9999999)
}

func TestZkNodeWatcher_Node_PERSISTENT_07(t *testing.T) {

	o := Config{
		Name:           "cc",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 3,
	}

	c, err := New(&o)
	assert.Equal(t, nil, err, "初始化失败")

	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}

}

func TestWatchNodeData(t *testing.T) {
	c, err := newZookeeperProxy(&Config{
		Name:           "cc",
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 3,
	})
	assert.Equal(t, nil, err, "初始化失败")

	path := "/mytest0006"
	c.CreateNode(path, PERSISTENT, "hello world")
	c.WatchNode(path, func(respChan <-chan *WatchNodeResponse) {
		for d := range respChan {
			fmt.Println("正常监听输出 "+path+" "+d.NodeChangeInfo.OperateType.String()+"  "+d.NodeChangeInfo.OldData+"  "+d.NodeChangeInfo.NewData, "  ", d.Err)
		}
	})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.SetData(path, "change-"+strconv.Itoa(i))
			fmt.Println(path, "change-"+strconv.Itoa(i))
		}(i)
	}
	wg.Wait()
	time.Sleep(3 * time.Second)
	c.Delete(path)
	c.Close()

}
