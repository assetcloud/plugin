#!/bin/sh
# proto生成命令，将pb.go文件生成到types/目录下, chain_path支持引用chain框架的proto文件
chain_path=$(go list -f '{{.Dir}}' "github.com/assetcloud/chain")
protoc --go_out=plugins=grpc:../../../../../../../ ./*.proto --proto_path=. --proto_path="${chain_path}/types/proto/"
