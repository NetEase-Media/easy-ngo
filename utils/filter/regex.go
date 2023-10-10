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
	"regexp"
	"strings"
)

type RegexOptions struct {
	Include *RegexOptionMeta
	Exclude *RegexOptionMeta
}

type RegexOptionMeta struct {
	Prefix []string
	Regex  []string
}

type RegexFilter struct {
	Rule *RegexOptions
}

func NewRegexFilter(opts *RegexOptions) *RegexFilter {
	if opts == nil {
		opts = &RegexOptions{
			Include: &RegexOptionMeta{},
			Exclude: &RegexOptionMeta{},
		}
	}
	if opts.Exclude == nil {
		opts.Exclude = &RegexOptionMeta{}
	}
	if opts.Include == nil {
		opts.Include = &RegexOptionMeta{}
	}

	return &RegexFilter{
		Rule: opts,
	}
}

func (f *RegexFilter) DoFilt(ctx context.Context, path string) (bool, error) {
	for _, prefix := range f.Rule.Include.Prefix {
		if strings.HasPrefix(path, prefix) {
			return false, nil
		}
	}
	for _, regex := range f.Rule.Include.Regex {
		re := regexp.MustCompile(regex)
		if re.MatchString(path) {
			return false, nil
		}
	}
	for _, prefix := range f.Rule.Exclude.Prefix {
		if strings.HasPrefix(path, prefix) {
			return true, nil
		}
	}
	for _, regex := range f.Rule.Exclude.Regex {
		re := regexp.MustCompile(regex)
		if re.MatchString(path) {
			return true, nil
		}
	}
	return false, nil // 都没有配置的则不需要过滤掉
}

func (f *RegexFilter) AddInclude(ctx context.Context, path string) error {
	f.Rule.Include.Regex = append(f.Rule.Include.Regex, path)
	return nil
}

func (f *RegexFilter) AddExclude(ctx context.Context, path string) error {
	f.Rule.Exclude.Regex = append(f.Rule.Exclude.Regex, path)
	return nil
}

func (f *RegexFilter) AddIncludePrefix(ctx context.Context, path string) error {
	f.Rule.Include.Prefix = append(f.Rule.Include.Prefix, path)
	return nil
}

func (f *RegexFilter) AddExcludePrefix(ctx context.Context, path string) error {
	f.Rule.Exclude.Prefix = append(f.Rule.Exclude.Prefix, path)
	return nil
}
