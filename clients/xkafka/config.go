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
	opt.DialTimeout = time.Second * 10
	opt.ReadTimeout = time.Second * 10
	opt.WriteTimeout = time.Second * 10
	opt.Metadata.Retries = 3
	opt.Metadata.Timeout = time.Second * 10
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
