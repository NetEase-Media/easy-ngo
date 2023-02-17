module github.com/NetEase-Media/easy-ngo/application

go 1.18

require (
	github.com/NetEase-Media/easy-ngo/clients/httplib v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/clients/xgorm v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/clients/xkafka v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/clients/xmemcache v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/clients/xredis v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/clients/xxxljob v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/clients/xzk v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/config v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/config/contrib/xagollo v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/microservices v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/microservices/contrib/sd/etcd v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/observability v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/observability/contrib/xotel v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/observability/contrib/xprometheus v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/servers/xgin v0.0.0-20230217025809-cc6ac801b1e8
	github.com/NetEase-Media/easy-ngo/xlog v0.0.0-20230208101755-f84181b2cdac
	github.com/NetEase-Media/easy-ngo/xlog/contrib/xzap v0.0.0-20230208101755-f84181b2cdac
	github.com/fatih/color v1.13.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
	go.etcd.io/etcd/client/v3 v3.5.6
	go.uber.org/multierr v1.7.0
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/NetEase-Media/easy-ngo/clients/xsentinel v0.0.0-20230208101755-f84181b2cdac // indirect
	github.com/Shopify/sarama v1.37.2 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/alibaba/sentinel-golang v1.0.4 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20221031212613-62deef7fc822 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/djimenez/iconv-go v0.0.0-20160305225143-8960e66bd3da // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-basic/ipv4 v1.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-redis/redis/extra/rediscmd v0.2.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/philchia/agollo v2.1.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/shirou/gopsutil/v3 v3.21.6 // indirect
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.43.0 // indirect
	github.com/xxl-job/xxl-job-executor-go v1.1.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.6 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.11.2 // indirect
	go.opentelemetry.io/otel/sdk v1.11.2 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20221207170731-23e4bf6bdc37 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.4.4 // indirect
	gorm.io/gorm v1.24.2 // indirect
)
