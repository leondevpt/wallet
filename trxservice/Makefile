GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

BUILDTIME=$(shell TZ=Asia/Shanghai date +%FT%T%z)
GitTag=$(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
GitBranch=$(shell git rev-parse --abbrev-ref HEAD)
GitCommit=$(shell git rev-parse --short=12 HEAD)
GitTreeState=$(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

REPO = trxservice

BUILD_FLAGS := -ldflags "-X '${REPO}/version.GitCommit=$(GitCommit)' \
                         -X '${REPO}/version.GitBranch=$(GitBranch)' \
                         -X '${REPO}/version.GitTag=$(GitTag)' \
                         -X '${REPO}/version.GitTreeState=$(GitTreeState)' \
                         -X '${REPO}/version.Version=$(VERSION)' \
                         -X '${REPO}/version.BuildTime=[$(BUILDTIME)]'"

LDFLAGS := ' -w -s'

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api/v1 \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api/v1 \
 	       --go-grpc_out=paths=source_relative:./api/v1 \
 	       --openapi_out==paths=source_relative:./api/v1 \
		   --grpc-gateway_out=paths=source_relative:./api/v1 \
	       $(API_PROTO_FILES)




.PHONY: build
# build
build:
	echo $(BUILD_FLAGS)
	mkdir -p bin/ && go build $(BUILD_FLAGS)$(LDFLAGS) -v -o ./bin/ ./...
