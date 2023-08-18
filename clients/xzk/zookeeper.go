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
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-zookeeper/zk"
)

const (
	PERSISTENT            = iota // 持久化节点
	EPHEMERAL                    // 临时节点， 客户端session超时这类节点就会被自动删除
	PERSISTENT_SEQUENTIAL        // 顺序自动编号持久化节点，这种节点会根据当前已存在的节点数自动加 1
	EPHEMERAL_SEQUENTIAL         // 临时自动编号节点
)

const (
	EventRecoverEphemeralNodeFail    EventType = 1
	EventRecoverEphemeralNodeSuccess EventType = 2
	EventSession                     EventType = -1
	Unknown                          EventType = -2
)

func newZookeeperProxy(opt *Config) (*ZookeeperProxy, error) {
	conn, ch, err := zk.Connect(opt.Addr, opt.SessionTimeout)
	if err != nil {
		return nil, fmt.Errorf("connection failed")
	}

	var scene sync.Map
	zp := &ZookeeperProxy{
		Config:   opt,
		Conn:     conn,
		listenCh: nil,
		tmpNode:  scene,
		stop:     make(chan struct{}),
		done:     &sync.WaitGroup{},
	}
	zp.done.Add(1)
	go zp.netLi(ch)
	return zp, nil
}

// 获取zookeeper当前的连接状态
func (z *ZookeeperProxy) GetConnState() string {
	s := z.Conn.State()
	if name := s.String(); name != "" {
		return name
	}
	return "Unknown"
}

// 获取当前连接的sessionId
func (z *ZookeeperProxy) GetSessionId() int64 {
	return z.Conn.SessionID()
}

type NodeDetail struct {
	Data string
	Acl  []zk.ACL
}

type EventType int32

// 监听连接状态
func (z *ZookeeperProxy) StartLi() (<-chan Event, error) {
	if z.listenCh != nil {
		return nil, errors.New("status listening has started")
	}
	z.listenCh = make(chan Event, 100)
	z.openListenSign()
	return z.listenCh, nil
}

type Event struct {
	Type   EventType
	State  zk.State
	Path   string
	Data   string
	Err    error
	Server string
}

func (z *ZookeeperProxy) netLi(c <-chan zk.Event) {
	defer func() {
		if !z.isListenSignClosed() {
			z.closeListenSign()
		}
		if z.listenCh != nil {
			close(z.listenCh)
		}
		z.done.Done()
	}()

	first := true
	var callback Event
	for {
		select {
		case sessionState := <-c:
			if !z.isListenSignClosed() {
				if sessionState.State.String() == "Unknown" {
					continue
				}
				switch sessionState.Type {
				case zk.EventSession:
					callback.Type = EventSession
				default:
					callback.Type = Unknown
				}
				callback.State = sessionState.State
				callback.Path = sessionState.Path
				callback.Data = ""
				callback.Err = sessionState.Err
				callback.Server = sessionState.Server
				z.listenCh <- callback
			}
			if sessionState.State == zk.StateHasSession {
				if first {
					first = false
					continue
				}
				z.tmpNode.Range(func(k, v interface{}) bool {
					path, _ := k.(string)
					d, _ := z.tmpNode.Load(k)
					detail := d.(*NodeDetail)
					dataByte := []byte(detail.Data)
					acls := detail.Acl
					_, ee := z.Conn.Create(path, dataByte, EPHEMERAL, acls)
					if ee != nil {
						for i := 0; i < 10 && ee != nil; i++ {
							_, ee = z.Conn.Create(path, dataByte, EPHEMERAL, acls)
						}
					}
					if ee != nil {
						callback.Type = EventRecoverEphemeralNodeFail
						callback.State = sessionState.State
						callback.Path = path
						callback.Data = detail.Data
						callback.Err = ee
						callback.Server = sessionState.Server
					} else {
						callback.Type = EventRecoverEphemeralNodeSuccess
						callback.State = sessionState.State
						callback.Path = path
						callback.Data = detail.Data
						callback.Err = ee
						callback.Server = sessionState.Server
					}
					if !z.isListenSignClosed() {
						z.listenCh <- callback
					}
					return true
				})
			}
		case <-z.stop:
			return
		}
	}
}

// 创建节点(拥有所有权限)
func (z *ZookeeperProxy) CreateNode(path string, flag int32, data string) (string, error) {
	dataByte := []byte(data)
	acls := zk.WorldACL(zk.PermAll) // 获取访问控制权限,默认全部权限
	// flags有4种取值:
	//   0:永久,除非手动删除
	//   1:短暂,session断开则该节点也被删除
	//   2:永久，且会自动在节点后面添加序号
	//   3:短暂且自动添加序号
	path, err := z.Conn.Create(path, dataByte, flag, acls)
	if err != nil {
		return "", err
	} else {
		if flag == EPHEMERAL || flag == EPHEMERAL_SEQUENTIAL {
			nodeDetail := new(NodeDetail)
			nodeDetail.Data = data
			nodeDetail.Acl = acls
			z.tmpNode.Store(path, nodeDetail)
			z.WatchNode(path, func(respChan <-chan *WatchNodeResponse) {
				for d := range respChan {
					nodeDetail.Data = d.NodeChangeInfo.NewData
					z.tmpNode.Store(path, nodeDetail)
				}
			})
		}
		return path, nil
	}
}

// 创建节点(自定义节点权限)
func (z *ZookeeperProxy) CreateNodeWithAcls(path string, flag int32, acls []zk.ACL, data string) (string, error) {
	dataByte := []byte(data)
	path, err := z.Conn.Create(path, dataByte, flag, acls)
	if err != nil {
		return "", err
	} else {
		if flag == EPHEMERAL || flag == EPHEMERAL_SEQUENTIAL {
			nodeDetail := new(NodeDetail)
			nodeDetail.Data = data
			nodeDetail.Acl = acls
			z.tmpNode.Store(path, nodeDetail)
			z.WatchNode(path, func(respChan <-chan *WatchNodeResponse) {
				for d := range respChan {
					nodeDetail.Data = d.NodeChangeInfo.NewData
					z.tmpNode.Store(path, nodeDetail)
				}
			})
		}
		return path, nil
	}
}

// 创建节点（包含父节点): 若要创建节点的父节点不存在，则把 父节点 和 需求节点都创建出来。  注： 所有新建的父节点的flags 为0 ， 需求节点的flags为传入的 flags
func (z *ZookeeperProxy) CreateNodeParent(path string, flags int32, data string) (string, error) {
	sign, _, err := z.Conn.Exists(path)
	if err != nil || sign {
		return "", errors.New("this node has existed")
	}
	nodes := strings.Split(path, "/")
	partialNode := ""
	for k, _ := range nodes {
		if nodes[k] == "" {
			continue
		}
		partialNode = partialNode + "/" + nodes[k]
		if exist, err := z.Exist(partialNode); err == nil && !exist {
			var err error
			if k == len(nodes)-1 {
				partialNode, err = z.CreateNode(partialNode, flags, data)
			} else {
				partialNode, err = z.CreateNode(partialNode, PERSISTENT, "")
			}
			if err != nil {
				return "", err
			}
		}
	}
	return partialNode, nil
}

// 判断节点是否存在
func (z *ZookeeperProxy) Exist(path string) (bool, error) {
	sign, _, err := z.Conn.Exists(path)
	return sign, err
}

// 设置节点值
func (z *ZookeeperProxy) SetData(path string, s string) error {
	data := []byte(s)
	_, stat, err := z.Conn.Get(path)
	if err != nil {
		return err
	}
	_, err = z.Conn.Set(path, data, stat.Version)
	return err
}

// 获取节点值
func (z *ZookeeperProxy) GetData(path string) (string, error) {
	data, _, err := z.Conn.Get(path)
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// 获取子节点列表
func (z *ZookeeperProxy) GetChildren(path string) ([]string, error) {
	children, _, err := z.Conn.Children(path)
	if err != nil {
		return []string{}, err
	} else {
		childrenPaths := []string{}
		for _, p := range children {
			childrenPaths = append(childrenPaths, path+"/"+p)
		}
		return childrenPaths, nil
	}
}

// 删除节点
func (z *ZookeeperProxy) Delete(path string) error {
	_, stat, err := z.Conn.Get(path)
	if err != nil {
		return err
	}
	err_ := z.Conn.Delete(path, stat.Version)
	return err_
}

type EventChildrenNodeChangeType int32

const (
	EventChildrenNodeDelete EventChildrenNodeChangeType = iota
	EventChildrenNodeIncrease
)

var (
	EventChildrenNodeChangeNames = map[EventChildrenNodeChangeType]string{
		EventChildrenNodeDelete:   "EventChildrenNodeDelete",
		EventChildrenNodeIncrease: "EventChildrenNodeIncrease",
	}
)

func (e EventChildrenNodeChangeType) String() string {
	if name := EventChildrenNodeChangeNames[e]; name != "" {
		return name
	}
	return "Unknown"
}

type ChildrenChange struct {
	OperateType EventChildrenNodeChangeType // 监控到子节点的操作类型: 目前只支持监控子节点的 增加 删除
	Path        string                      // 节点路径
}
type WatchChildrenResponse struct {
	ChildrenChangeInfo []*ChildrenChange
	Err                error
}

// 监听子节点                               				TODO： 目前只支持对子节点进行 增加 或 删除的监控， 不支持对 子节点值改变 的监控
func (z *ZookeeperProxy) WatchChildren(path string, listener func(respChan <-chan *WatchChildrenResponse)) {
	z.done.Add(1)
	respChan := make(chan *WatchChildrenResponse, 100)
	go listener(respChan)

	oldChildren, _, eventCh, err := z.Conn.ChildrenW(path)
	if err != nil {
		errResp := &WatchChildrenResponse{
			ChildrenChangeInfo: []*ChildrenChange{},
			Err:                err,
		}
		respChan <- errResp
		return
	}
	go func(path string, oldChildren []string) {
		defer func() {
			close(respChan)
			z.done.Done()
		}()

		var response *WatchChildrenResponse
		for {
			select {
			case <-z.stop:
				return
			case <-eventCh:
				newChildren, _, newEventCh, er := z.Conn.ChildrenW(path)
				eventCh = newEventCh
				if er != nil {
					response = &WatchChildrenResponse{
						ChildrenChangeInfo: nil,
						Err:                er,
					}
					respChan <- response
					if er == zk.ErrNoNode {
						return
					}
				} else {
					response = childrenWatcherResponse(path, oldChildren, newChildren)
					respChan <- response
					oldChildren = newChildren
				}
			}
		}
	}(path, oldChildren)
}

func childrenWatcherResponse(path string, oldChildren, newChildren []string) *WatchChildrenResponse {
	wCResponse := &WatchChildrenResponse{}
	changeS := []*ChildrenChange{}
	mOldChildren, mNewChildren := make(map[string]int), make(map[string]int)
	for _, v := range oldChildren {
		mOldChildren[v] = 0
	}
	for _, vv := range newChildren {
		mNewChildren[vv] = 0
	}
	for k, _ := range mNewChildren {
		if _, exist := mOldChildren[k]; exist {
			delete(mOldChildren, k)
		} else {
			resp := &ChildrenChange{
				OperateType: EventChildrenNodeIncrease,
				Path:        path + "/" + k,
			}
			changeS = append(changeS, resp)
		}
	}

	for k, _ := range mOldChildren {
		resp := &ChildrenChange{
			OperateType: EventChildrenNodeDelete,
			Path:        path + "/" + k,
		}
		changeS = append(changeS, resp)
	}
	wCResponse.ChildrenChangeInfo = changeS
	wCResponse.Err = nil
	return wCResponse
}

type EventCurrentNodeChangeType int32

const (
	EventNodeDelete EventCurrentNodeChangeType = iota
	EventNodeDataChange
)

var (
	EventCurrntNodeChangeNames = map[EventCurrentNodeChangeType]string{
		EventNodeDelete:     "EventNodeDelete",
		EventNodeDataChange: "EventNodeDataChange",
	}
)

func (t EventCurrentNodeChangeType) String() string {
	if name := EventCurrntNodeChangeNames[t]; name != "" {
		return name
	}
	return "Unknown"
}

type NodeChange struct {
	OperateType EventCurrentNodeChangeType // 监控到的子节点的操作类型
	OldData     string                     // 节点改变前的值
	NewData     string                     // 节点改变后的值
	Path        string                     // 节点路径
}
type WatchNodeResponse struct {
	NodeChangeInfo *NodeChange
	Err            error
}

// 对节点监听
func (z *ZookeeperProxy) WatchNode(path string, listener func(respChan <-chan *WatchNodeResponse)) {
	z.done.Add(1)
	respChan := make(chan *WatchNodeResponse, 100)
	go listener(respChan)
	go func(path string) {
		defer func() {
			close(respChan)
			z.done.Done()
		}()

		var oldData string
		var response *WatchNodeResponse
		oldDataBuf, _, eventCh, err := z.Conn.GetW(path)
		if err != nil {
			response = &WatchNodeResponse{
				NodeChangeInfo: &NodeChange{},
				Err:            err,
			}
			respChan <- response
			return
		}
		_, isEphemeralNode := z.tmpNode.Load(path)
		oldData = string(oldDataBuf)
		for {
			select {
			case <-z.stop:
				return
			case e := <-eventCh:
				newDataBuf, _, newEventCh, er := z.Conn.GetW(path)
				for isEphemeralNode && er != nil {
					newDataBuf, _, newEventCh, er = z.Conn.GetW(path)
				}
				eventCh = newEventCh
				if er != nil && !isEphemeralNode && e.Type == zk.EventNodeDeleted && er == zk.ErrNoNode {
					response = &WatchNodeResponse{
						NodeChangeInfo: &NodeChange{OperateType: EventNodeDelete},
						Err:            er,
					}
					respChan <- response
					return
				} else if er == nil {
					if oldData == string(newDataBuf) {
						continue
					}
					nodeChange := &NodeChange{
						OperateType: EventNodeDataChange,
						OldData:     oldData,
						NewData:     string(newDataBuf),
						Path:        path,
					}
					response = &WatchNodeResponse{
						NodeChangeInfo: nodeChange,
						Err:            nil,
					}
					respChan <- response
				}
				oldData = string(newDataBuf)
			}
		}
	}(path)
}

func (z *ZookeeperProxy) Close() error {
	close(z.stop)
	z.done.Wait()
	z.Conn.Close()
	return nil
}

func (z *ZookeeperProxy) openListenSign() error {
	if !atomic.CompareAndSwapUint32(&z.listenSign, 0, 1) {
		return io.ErrClosedPipe
	}
	return nil
}

func (z *ZookeeperProxy) closeListenSign() error {
	if !atomic.CompareAndSwapUint32(&z.listenSign, 1, 0) {
		return io.ErrClosedPipe
	}
	return nil
}

func (z *ZookeeperProxy) isListenSignClosed() bool {
	return atomic.LoadUint32(&z.listenSign) == 0
}
