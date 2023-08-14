package pluginxgin

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	return nil
}
