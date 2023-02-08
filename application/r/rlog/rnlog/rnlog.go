package rlog

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngocation/r/rlog"
	conf "github.com/NetEase-Media/easy-ngog"

	"github.com/NetEase-Media/easy-ngonlog"
)

const (
	key_nlog = rlog.Key_prefix + ".nlog"
)

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	nlogOpts := make([]nlog.Option, 0)
	if err := conf.Get(key_nlog, &nlogOpts); err != nil {
		panic("load nlog config failed.")
	}
	if len(nlogOpts) == 0 {
		panic("no nlog config!")
	}
	for _, nlogOpt := range nlogOpts {
		nlog, err := nlog.New(&nlogOpt)
		if err != nil {
			panic("init nlog failed.")
		}
		rlog.Set(nlogOpt.Name, nlog)
	}
	return nil
}
