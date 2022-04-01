PROTOBUF_PATH?=./proto
PROTO_SERVICES = $(patsubst ${PROTOBUF_PATH}/%.proto, %.proto, $(wildcard ${PROTOBUF_PATH}/*.proto))

.PHONY: all

all: clean build

gen:
	mkdir -p internal/pb
	protoc \
		--go_out=internal/pb \
		--go-grpc_out=internal/pb \
		--go-grpc_opt=module=github.com/devopshaven/gateway-auth-service/internal/pb \
		--proto_path=${PROTOBUF_PATH} \
		--go_opt=module=github.com/devopshaven/gateway-auth-service/internal/pb ${PROTO_SERVICES}
