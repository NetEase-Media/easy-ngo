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
	"github.com/gorilla/mux"
	"net/http"
)

// if a path match the include list and exclude list, we only use include list
type RouteMuxFilter interface {
	AddInclude(ctx context.Context, path string) error
	AddExclude(ctx context.Context, path string) error
	BelongToInclude(req *http.Request, match *mux.RouteMatch) bool
	BelongToExclude(req *http.Request, match *mux.RouteMatch) bool
}

type DefaultRouteFilter struct {
	Include *mux.Router
	Exclude *mux.Router
}

type RouteOptions struct {
	Include []string
	Exclude []string
}

func NewDefaultRouteFilter(opts *RouteOptions) *DefaultRouteFilter {
	f := &DefaultRouteFilter{
		Include: mux.NewRouter(),
		Exclude: mux.NewRouter(),
	}
	for _, path := range opts.Include {
		f.AddInclude(context.Background(), path)
	}
	for _, path := range opts.Exclude {
		f.AddExclude(context.Background(), path)
	}
	return f
}

func (f *DefaultRouteFilter) AddInclude(ctx context.Context, path string) error {
	f.Include.Path(path)
	return nil
}

func (f *DefaultRouteFilter) AddExclude(ctx context.Context, path string) error {
	f.Include.Path(path)
	return nil
}

func (f *DefaultRouteFilter) BelongToInclude(req *http.Request, match *mux.RouteMatch) bool {
	has1 := f.findInclude(req, match)
	if has1 {
		return has1
	}
	return false
}

func (f *DefaultRouteFilter) BelongToExclude(req *http.Request, match *mux.RouteMatch) bool {
	has1 := f.findExclude(req, match)
	if has1 {
		return has1
	}
	return false
}

func (f *DefaultRouteFilter) findInclude(req *http.Request, match *mux.RouteMatch) bool {
	hd := f.Include.Match(req, match)
	return hd
}

func (f *DefaultRouteFilter) findExclude(req *http.Request, match *mux.RouteMatch) bool {
	hd := f.Exclude.Match(req, match)
	return hd
}
