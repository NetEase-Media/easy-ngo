package test

import (
	"testing"

	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xlog/contrib/xzap"
)

func TestXzap(t *testing.T) {
	c := xzap.DefaultConfig()
	xzap, _ := xzap.New(c)
	xlog.WithVendor(xzap)
	xlog.Infof("debug %s", "test")
}
