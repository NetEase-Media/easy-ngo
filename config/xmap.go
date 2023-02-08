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

package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	keyDelimiter = "."
)

var (
	ExpanderPattern = regexp.MustCompile("\\$\\{[^\\$]+\\}")
)

// 带有层级结构的mpa转换成 拍平的结构
// 例如：
// A
//
//	b
//	c
//
// 转变为
// A.b
// A.c
// 叶子的节点的格式都有哪些？
// 数组？ map？ 还有一些其他类型的值吗？
func Flattening(val map[string]interface{}) map[string]interface{} {
	if val == nil {
		return nil
	}
	M := make(map[string]interface{})
	// 根节点
	for k, v := range val {
		m := make(map[string]interface{})
		traverse("", k, v, m)
		Merge(M, m)
	}
	return M
}

func traverse(prefix, key string, value interface{}, target map[string]interface{}) {
	// 递归到不是 map 结构, 也就是到了叶子接点
	// 也可能结点是个数组，但是数组内的元素是个结构体，还是需要往下遍历的
	// 但是这种情况好像也不需要往下遍历了
	kk := fmt.Sprintf("%s%s%s", prefix, keyDelimiter, key)
	if prefix == "" { // root node, no need append delimiter
		kk = key
	}
	target[kk] = value
	switch value.(type) {
	case map[string]interface{}: // we only handle this type so far.
		for k, v := range value.(map[string]interface{}) {
			traverse(kk, k, v, target)
		}
	default:
		target[kk] = value
	}

}

// the parameter override can be written by source
// if the override has the same key with source.
func Merge(override map[string]interface{}, source map[string]interface{}) {
	if override == nil {
		override = source
	}
	if source == nil {
		return
	}
	for k, v := range source {
		override[k] = v // traverse and override.
	}
}

func Expand(original map[string]interface{}) (result map[string]interface{}) {
	if original == nil {
		return nil
	}
	result = make(map[string]interface{}, 0)
	for k, v := range original {
		result[k] = expandMap(original, "", k, v)
		//fmt.Printf("---result:%s", result)
	}
	return result
}

func expandMap(original map[string]interface{}, prefix string, key string, val interface{}) (ret interface{}) {
	prefix = prefix + key + keyDelimiter

	//fmt.Printf("expandMap:%s\n", prefix)

	switch val.(type) {
	case map[string]interface{}: // we only handle this type so far.
		retMap := make(map[string]interface{})
		for k, v := range val.(map[string]interface{}) {
			retMap[k] = expandMap(original, prefix, k, v)
		}
		ret = retMap
	case []map[string]interface{}:
		vals := val.([]map[string]interface{})
		retSlice := make([]map[string]interface{}, 0, len(vals))
		for i, v := range vals {
			temp := expandMap(original, prefix, strconv.FormatInt(int64(i), 10), v)
			retSlice = append(retSlice, temp.(map[string]interface{}))
		}
		ret = retSlice
	case string:
		valString := val.(string)
		ret = expandString(original, prefix, key, valString)
	default:
		ret = val
	}
	return ret
}

func expandString(original map[string]interface{}, prefix string, key string, val string) (ret string) {
	//fast test?
	if strings.Contains(val, "${") {
		return ExpanderPattern.ReplaceAllStringFunc(val, func(match string) string {
			// 获取到引用的key
			inner := match[2 : len(match)-1]
			// 按分隔符切分
			keys := strings.Split(inner, keyDelimiter)
			// 按原始数据展开
			result, err := expandStringByKeys(original, keys)
			// 遇到错误使用原值
			if err != nil {
				fmt.Printf("expr:%s error.\n", match)
				return match
			}
			return result
		})
	}
	return val
}

// 按原始数据展开,只支持叶节点的引用,仅支持在string中引用
// 支持引用的数据类型有int,bool,float,string
// 其他方式的引用不做处理
func expandStringByKeys(original map[string]interface{}, keys []string) (string, error) {
	//fmt.Printf("keys:%s", keys)
	var tempI interface{} = original
	for _, key := range keys {
		switch tempI.(type) {
		case map[string]interface{}:
			tempI = tempI.(map[string]interface{})[key]
		case []map[string]interface{}:
			keyIndex, _ := strconv.ParseInt(key, 10, 64)
			tempI = tempI.([]map[string]interface{})[keyIndex]
		default:
			fmt.Printf("key:%s,Unknown type:%s\n", key, reflect.TypeOf(tempI))
		}
	}
	// 如果引用的值不存在，则返回error
	if tempI == nil {
		return "", errors.New("nil")
	}
	ret := ""
	var err error
	switch tempI.(type) {
	case int, int8, int16, int32, int64:
		ret = strconv.FormatInt(tempI.(int64), 10)
	case bool:
		ret = strconv.FormatBool(tempI.(bool))
	case float32, float64:
		ret = strconv.FormatFloat(tempI.(float64), 'e', 10, 64)
	case string:
		ret = tempI.(string)
	//case time.Time:
	//	ret = tempI.(time.Time).Format("")
	default:
		err = errors.New("Unsupport type")
	}
	return ret, err
}
