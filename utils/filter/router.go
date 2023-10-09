// Copyright 2023 NetEase Media Technology（Beijing）Co., Ltd.
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

package filter

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// if a path match the include list and exclude list, we only use include list
type HttpFilter interface {
	AddInclude(ctx context.Context, path string) error
	AddExclude(ctx context.Context, path string) error
	DoFilt(ctx context.Context, path string) (bool, error)
}

// use httprouter implement HttpFilter
type DefaultFilter struct {
	Include *httprouter.Router
	Exclude *httprouter.Router
}

type Options struct {
	Include []string
	Exclude []string
}

func NewDefaultFilter(opts *Options) *DefaultFilter {
	f := &DefaultFilter{
		Include: httprouter.New(),
		Exclude: httprouter.New(),
	}
	for _, path := range opts.Include {
		f.AddInclude(context.Background(), path)
	}
	for _, path := range opts.Exclude {
		f.AddExclude(context.Background(), path)
	}
	return f
}

var (
	hand          = func(http.ResponseWriter, *http.Request, httprouter.Params) {}
	defaultMethod = "*" // think not support method
)

func (f *DefaultFilter) AddInclude(ctx context.Context, path string) error {
	f.Include.Handle(defaultMethod, path, hand)
	return nil
}

func (f *DefaultFilter) AddExclude(ctx context.Context, path string) error {
	f.Exclude.Handle(defaultMethod, path, hand)
	return nil
}

func (f *DefaultFilter) DoFilt(ctx context.Context, path string) (bool, error) {
	has1, err := f.findInclude(ctx, path)
	if err != nil {
		return false, err
	}
	if has1 { // 不需要被过滤掉，因为在include列表中
		return false, nil
	}
	has2, err := f.findExclude(ctx, path)
	if err != nil {
		return false, err
	}
	if has2 { // 需要被过滤掉，因为在exclude列表中
		return true, nil
	}
	return false, nil
}

func (f *DefaultFilter) findInclude(ctx context.Context, path string) (bool, error) {
	hd, _, _ := f.Include.Lookup(defaultMethod, path)
	if hd != nil {
		return true, nil
	}
	return false, nil
}

func (f *DefaultFilter) findExclude(ctx context.Context, path string) (bool, error) {
	hd, _, _ := f.Exclude.Lookup(defaultMethod, path)
	if hd != nil {
		return true, nil
	}
	return false, nil
}
