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

package xkafka

import (
	"time"

	"github.com/IBM/sarama"
)

type Config struct {
	Name            string
	Addr            []string
	Version         string
	MaxOpenRequests int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	SASL            struct {
		Enable                   bool
		Mechanism                sarama.SASLMechanism
		Version                  int16
		Handshake                bool
		AuthIdentity             string
		User                     string
		Password                 string
		SCRAMAuthzID             string
		SCRAMClientGeneratorFunc func() sarama.SCRAMClient
		TokenProvider            sarama.AccessTokenProvider
		GSSAPI                   sarama.GSSAPIConfig
	}
	Metadata struct {
		Retries int
		Timeout time.Duration
	}
	Consumer struct {
		Group              string
		EnableAutoCommit   bool
		AutoCommitInterval time.Duration
		InitialOffset      int64
		SessionTimeout     time.Duration
		MinFetchBytes      int32
		DefaultFetchBytes  int32
		MaxFetchBytes      int32
		MaxFetchWait       time.Duration
		Retries            int
	}
	Producer struct {
		MaxMessageBytes  int
		Acks             sarama.RequiredAcks
		Timeout          time.Duration
		Retries          int
		MaxFlushBytes    int
		MaxFlushMessages int
		FlushFrequency   time.Duration
		Idempotent       bool
	}
	EnableMetrics bool
	EnableTrace   bool
}

func DefaultConfig() *Config {
	opt := &Config{}
	opt.Version = defaultVersion
	opt.MaxOpenRequests = 5
	opt.DialTimeout = time.Second * 30
	opt.ReadTimeout = time.Second * 30
	opt.WriteTimeout = time.Second * 30
	opt.Metadata.Retries = 3
	opt.Metadata.Timeout = time.Second * 60
	opt.Consumer.Group = ""
	opt.Consumer.EnableAutoCommit = true
	opt.Consumer.AutoCommitInterval = time.Second * 1
	opt.Consumer.InitialOffset = sarama.OffsetNewest
	opt.Consumer.SessionTimeout = time.Second * 10
	opt.Consumer.MinFetchBytes = 1
	opt.Consumer.DefaultFetchBytes = 1024 * 1024
	opt.Consumer.MaxFetchBytes = 0
	opt.Consumer.MaxFetchWait = time.Millisecond * 250
	opt.Consumer.Retries = 3
	opt.Producer.MaxMessageBytes = 1000000
	opt.Producer.Acks = sarama.WaitForLocal
	opt.Producer.Timeout = time.Second * 10
	opt.Producer.Retries = 3
	opt.Producer.MaxFlushBytes = 100 * 1024 * 1024
	opt.Producer.MaxFlushMessages = 0
	opt.Producer.FlushFrequency = time.Second * 1
	opt.Producer.Idempotent = false
	return opt
}
