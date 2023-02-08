package xkafka

import (
	"errors"
	"time"

	"github.com/Shopify/sarama"
)

type Option struct {
	Name            string
	Addr            []string
	Version         string
	MaxOpenRequests int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	SASL            struct {
		// Whether or not to use SASL authentication when connecting to the broker
		// (defaults to false).
		Enable bool
		// SASLMechanism is the name of the enabled SASL mechanism.
		// Possible values: OAUTHBEARER, PLAIN (defaults to PLAIN).
		Mechanism sarama.SASLMechanism
		// Version is the SASL Protocol Version to use
		// Kafka > 1.x should use V1, except on Azure EventHub which use V0
		Version int16
		// Whether or not to send the Kafka SASL handshake first if enabled
		// (defaults to true). You should only set this to false if you're using
		// a non-Kafka SASL proxy.
		Handshake bool
		// AuthIdentity is an (optional) authorization identity (authzid) to
		// use for SASL/PLAIN authentication (if different from User) when
		// an authenticated user is permitted to act as the presented
		// alternative user. See RFC4616 for details.
		AuthIdentity string
		// User is the authentication identity (authcid) to present for
		// SASL/PLAIN or SASL/SCRAM authentication
		User string
		// Password for SASL/PLAIN authentication
		Password string
		// authz id used for SASL/SCRAM authentication
		SCRAMAuthzID string
		// SCRAMClientGeneratorFunc is a generator of a user provided implementation of a SCRAM
		// client used to perform the SCRAM exchange with the server.
		SCRAMClientGeneratorFunc func() sarama.SCRAMClient
		// TokenProvider is a user-defined callback for generating
		// access tokens for SASL/OAUTHBEARER auth. See the
		// AccessTokenProvider interface docs for proper implementation
		// guidelines.
		TokenProvider sarama.AccessTokenProvider

		GSSAPI sarama.GSSAPIConfig
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
	EnableTracer  bool
}

func NewDefaultOptions() *Option {
	opt := &Option{}
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

func checkOptions(opt *Option) error {
	if opt.Name == "" {
		return errors.New("client name can not be nil")
	}
	if opt.Version == "" {
		opt.Version = defaultVersion
	}

	if len(opt.Addr) == 0 {
		return errors.New("empty address")
	}
	return nil
}

func (opt *Option) fulfill() {
	if opt.MaxOpenRequests == 0 {
		opt.MaxOpenRequests = 5
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = time.Second * 30
	}
	if opt.ReadTimeout == 0 {
		opt.ReadTimeout = time.Second * 30
	}
	if opt.WriteTimeout == 0 {
		opt.WriteTimeout = time.Second * 30
	}
	if opt.Producer.Timeout == 0 {
		opt.Producer.Timeout = time.Second * 10
	}
	/*
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
		opt.Producer.Retries = 3
		opt.Producer.MaxFlushBytes = 100 * 1024 * 1024
		opt.Producer.MaxFlushMessages = 0
		opt.Producer.FlushFrequency = time.Second * 1
		opt.Producer.Idempotent = false
	*/
}
