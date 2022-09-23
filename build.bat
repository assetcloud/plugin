go env -w CGO_ENABLED=0
go build -o chain.exe
go build -o chain-cli.exe github.com/assetcloud/plugin/cli
