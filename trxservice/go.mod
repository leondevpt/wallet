module github.com/leondevpt/wallet/trxservice

go 1.18

replace github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.22.1

require (
	github.com/apex/log v1.9.0
	github.com/fbsobreira/gotron-sdk v0.0.0-20211102183839-58a64f4da5f4
	github.com/fsnotify/fsnotify v1.5.4
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/joho/godotenv v1.4.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.13.0
	github.com/spf13/viper v1.7.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.opencensus.io v0.22.5
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.1
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	google.golang.org/grpc v1.50.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.3.6
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/extra/rediscmd v0.2.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/ethereum/go-ethereum v1.10.25 // indirect
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/protobuf v1.5.2
	github.com/google/wire v0.5.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/shengdoushi/base58 v1.0.0 // indirect
	github.com/shopspring/decimal v1.3.1
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20221004154528-8021a29435af
	golang.org/x/sys v0.0.0-20221006211917-84dc82d7e875 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220930163606-c98284e70a91
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/gorm v1.23.10
	moul.io/zapgorm2 v1.1.3
)
