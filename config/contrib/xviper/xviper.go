package xviper

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type XViper struct {
	viper *viper.Viper
}

func New() *XViper {
	return &XViper{
		viper: viper.New(),
	}
}

func (xviper *XViper) Init(protocol string) error {
	scheme := protocol[:strings.Index(protocol, "://")]
	switch scheme {
	case "env":
		initEnv(protocol[strings.Index(protocol, "://")+3:])
	case "file":
		initFile(protocol[strings.Index(protocol, "://")+3:])
	}
	return nil
}

func initEnv(protocol string) {
	kvs := strings.Split(protocol, ";")
	for _, kv := range kvs {
		if strings.HasPrefix(kv, "prefix=") {
			xviper.viper.SetEnvPrefix(strings.TrimPrefix(kv, "prefix="))
		}
	}
	xviper.viper.AutomaticEnv()
}

func initFile(protocol string) {
	kvs := strings.Split(protocol, ";")
	for _, kv := range kvs {
		if strings.HasPrefix(kv, "name=") {
			xviper.viper.SetConfigName(strings.TrimPrefix(kv, "name="))
		} else if strings.HasPrefix(kv, "type=") {
			xviper.viper.SetConfigType(strings.TrimPrefix(kv, "type="))
		} else if strings.HasPrefix(kv, "path=") {
			paths := strings.Split(strings.TrimPrefix(kv, "path="), ",")
			for _, path := range paths {
				xviper.viper.AddConfigPath(path)
			}
		}
	}
}

func (xviper *XViper) Read() error {
	return xviper.viper.ReadInConfig()
}

func (xviper *XViper) GetString(key string) string {
	return xviper.viper.GetString(key)
}

func (xviper *XViper) UnmarshalKey(key string, rawVal interface{}) error {
	return xviper.viper.UnmarshalKey(key, rawVal)
}

func (xviper *XViper) GetInt(key string) int {
	return xviper.viper.GetInt(key)
}

func (xviper *XViper) GetBool(key string) bool {
	return xviper.viper.GetBool(key)
}

func (xviper *XViper) GetTime(key string) time.Time {
	return xviper.viper.GetTime(key)
}

func (xviper *XViper) GetFloat64(key string) float64 {
	return xviper.viper.GetFloat64(key)
}
