
# golang1.12 or latest
# 1. make help
# 2. make dep
# 3. make build
# ...
export GO111MODULE=on
CLI := build/chain-cli
SRC_CLI := github.com/assetcloud/plugin/cli
APP := build/chain
LDFLAGS := ' -w -s'
BUILDTIME:=$(shell date +"%Y-%m-%d %H:%M:%S %A")
VERSION=$(shell git describe --tags || git rev-parse --short=8 HEAD)
GitCommit=$(shell git rev-parse --short=8 HEAD)
BUILD_FLAGS := -ldflags '-X "github.com/assetcloud/plugin/version.GitCommit=$(GitCommit)" \
                         -X "github.com/assetcloud/plugin/version.Version=$(VERSION)" \
                         -X "github.com/assetcloud/plugin/version.BuildTime=$(BUILDTIME)"'

MKPATH=$(abspath $(lastword $(MAKEFILE_LIST)))
MKDIR=$(dir $(MKPATH))
proj := "build"
.PHONY: default dep all build release cli linter race test fmt vet bench msan coverage coverhtml docker docker-compose protobuf clean help autotest pkg

default: depends build

build: depends
	go build $(BUILD_FLAGS) -v -o $(APP)
	go build $(BUILD_FLAGS) -v -o $(CLI) $(SRC_CLI)
	go build $(BUILD_FLAGS) -v -o build/fork-config github.com/assetcloud/plugin/cli/fork_config/
	@cp chain.toml build/
	@cp chain.para.toml build/ci/paracross/


build_ci: depends ## Build the binary file for CI
	@go build -v -o $(CLI) $(SRC_CLI)
	@go build $(BUILD_FLAGS) -v -o $(APP)
	@cp chain.toml build/
	@cp chain.para.toml build/ci/paracross/

pkg:
	rm assetchain-qbft assetchain-qbft.tgz -rf
	mkdir assetchain-qbft
	cp build/chain build/chain-cli chain.para.toml genconfig.sh assetchain-qbft
	tar zcfv assetchain-qbft.tgz assetchain-qbft

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \

WINDOWS_ARCH_LIST = \
	windows-amd64

GOBUILD=go build $(BUILD_FLAGS)" -w -s"

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(APP)-$@ $(SRC)
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(CLI)-$@ $(SRC_CLI)
	cp chain.para.toml chain.toml CHANGELOG.md build/ && cd build && \
	chmod +x chain-darwin-amd64 && \
	chmod +x chain-cli-darwin-amd64 && \
	tar -zcvf chain-darwin-amd64.tar.gz chain-darwin-amd64 chain-cli-darwin-amd64 chain.para.toml chain.toml CHANGELOG.md

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(APP)-$@ $(SRC)
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(CLI)-$@ $(SRC_CLI)
	cp chain.para.toml chain.toml CHANGELOG.md build/ && cd build && \
	chmod +x chain-darwin-arm64 && \
	chmod +x chain-cli-darwin-arm64 && \
	tar -zcvf chain-darwin-arm64.tar.gz chain-darwin-arm64 chain-cli-darwin-arm64 chain.toml chain.para.toml CHANGELOG.md

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(APP)-$@ $(SRC)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(CLI)-$@ $(SRC_CLI)
	cp chain.para.toml chain.toml CHANGELOG.md build/ && cd build && \
	chmod +x chain-linux-amd64 && \
	chmod +x chain-cli-linux-amd64 && \
	tar -zcvf chain-linux-amd64.tar.gz chain-linux-amd64 chain-cli-linux-amd64 chain.para.toml chain.toml CHANGELOG.md

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(APP)-$@.exe $(SRC)
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(CLI)-$@.exe $(SRC_CLI)
	cp chain.para.toml chain.toml CHANGELOG.md build/ && cd build && \
	zip -j  chain-windows-amd64.zip chain-windows-amd64.exe chain-cli-windows-amd64.exe chain.para.toml chain.toml CHANGELOG.md

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

para:
	@go build -v -o build/$(NAME) -ldflags "-X $(SRC_CLI)/buildflags.ParaName=user.p.$(NAME). -X $(SRC_CLI)/buildflags.RPCAddr=http://localhost:8901" $(SRC_CLI)

vet:
	@go vet -copylocks=false ./...

autotest: ## build autotest binary
	@cd build/autotest && bash ./run.sh build && cd ../../
	@if [ -n "$(dapp)" ]; then \
	cd build/autotest && bash ./run.sh local $(dapp) && cd ../../; fi
autotest_ci: autotest ## autotest ci
	@cd build/autotest && bash ./run.sh jerkinsci $(proj) && cd ../../
autotest_tick: autotest ## run with ticket mining
	@cd build/autotest && bash ./run.sh gitlabci build && cd ../../

update: ## version 可以是git tag打的具体版本号,也可以是commit hash, 什么都不填的情况下默认从master分支拉取最新版本
	@if [ -n "$(version)" ]; then   \
	go get -v github.com/assetcloud/chain@${version}  ; \
	else \
	go get -v github.com/assetcloud/chain@master ;fi
	@go mod tidy
dep:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.18.0
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/mitchellh/gox
	@go get -u github.com/vektra/mockery/.../
	@go get -u mvdan.cc/sh/cmd/shfmt
	@go get -u mvdan.cc/sh/cmd/gosh
	@git checkout go.mod go.sum
	@apt install clang-format
	@apt install shellcheck

linter: vet ineffassign ## Use gometalinter check code, ignore some unserious warning
	@./golinter.sh "filter"
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shellcheck -e SC2086

linter_test: ## Use gometalinter check code, for local test
	@./golinter.sh "test" "${p}"
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shellcheck -e SC2294

ineffassign:
	@golangci-lint  run --no-config --issues-exit-code=1  --deadline=2m --disable-all   --enable=ineffassign ./...

race: ## Run data race detector
	@go test -parallel=8 -race -short `go list ./... | grep -v "pbft"`

test: ## Run unittests
	@go test -parallel=8 -race  `go list ./...| grep -v "pbft"`

testq: ## Run unittests
	@go test -parallel=8 `go list ./... | grep -v "pbft"`

fmt: fmt_proto fmt_shell ## go fmt
	@go fmt ./...
	@find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w

.PHONY: fmt_proto fmt_shell
fmt_proto: ## go fmt protobuf file
	#@find . -name '*.proto' -not -path "./vendor/*" | xargs clang-format -i

fmt_shell: ## check shell file
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shfmt -w -s -i 4 -ci -bn

fmt_go: fmt_shell ## go fmt
	@go fmt ./...
	@find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w


coverage: ## Generate global code coverage report
	@./build/tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	@./build/tools/coverage.sh html;

docker: ## build docker image for chain run
	@sudo docker build . -f ./build/Dockerfile-run -t chain:latest

#extra can make more test setting
docker-compose: ## build docker-compose for chain run
	@cd build && if ! [ -d ci ]; then \
	 make -C ../ ; \
	 fi; \
	 cp chain* Dockerfile  docker-compose.yml *.sh ci/ && cd ci/ && ./docker-compose-pre.sh run $(proj) $(dapp) $(extra) && cd ../..

docker-compose-down: ## build docker-compose for chain run
	@cd build && if [ -d ci ]; then \
	 cp chain* Dockerfile  docker-compose* ci/ && cd ci/ && ./docker-compose-pre.sh down $(proj) $(dapp) && cd .. ; \
	 fi; \
	 cd ..

metrics:## build docker-compose for chain metrics
	@cd build && if ! [ -d ci ]; then \
	 make -C ../ ; \
	 fi; \
	 cp chain* Dockerfile  docker-compose.yml docker-compose-metrics.yml influxdb.conf *.sh ci/paracross/testcase.sh metrics/ && ./docker-compose-pre.sh run $(proj) metrics  && cd ../..


fork-test: ## build fork-test for chain run
	@cd build && cp chain* Dockerfile system-fork-test.sh docker-compose* ci/ && cd ci/ && ./docker-compose-pre.sh forktest $(proj) $(dapp) && cd ../..

largefile-check:
	git gc
	./findlargefile.sh

clean: ## Remove previous build
	@rm -rf $(shell find . -name 'datadir' -not -path "./vendor/*")
	@rm -rf build/chain*
	@rm -rf build/relayd*
	@rm -rf build/*.log
	@rm -rf build/logs
	@rm -rf build/autotest/autotest
	@rm -rf build/ci
	@rm -rf build/system-rpc-test.sh
	@rm -rf tool
	@cd build/metrics && find * -not -name readme.md | xargs rm -fr && cd ../..
	@go clean

proto:protobuf

protobuf: ## Generate protbuf file of types package
#	@cd ${CHAIN_PATH}/types/proto && ./create_protobuf.sh && cd ../..
	@find ./plugin/dapp -maxdepth 2 -type d  -name proto -exec make -C {} \;

depends: ## Generate depends file of types package
	@find ./plugin/dapp -maxdepth 2 -print -type d  -name cmd -exec make -C {} OUT="$(MKDIR)build/ci" FLAG= \;


help: ## Display this help screen
	@printf "Help doc:\nUsage: make [command]\n"
	@printf "[command]\n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cleandata:
	rm -rf build/datadir/addrbook
	rm -rf build/datadir/blockchain.db
	rm -rf build/datadir/mavltree
	rm -rf build/chain.log

.PHONY: checkgofmt
checkgofmt: ## get all go files and run go fmt on them
	@files=$$(find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -l -s); if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  echo "${files}"; \
		  exit 1; \
		  fi;
	@files=$$(find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w); if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  echo "${files}"; \
		  exit 1; \
		  fi;

.PHONY: auto_ci_before auto_ci_after auto_ci
auto_ci_before: clean fmt protobuf
	@echo "auto_ci"
	@go version
	@protoc --version
	@mockery -version
	@docker version
	@docker-compose version
	@git version
	@git status

.PHONY: auto_ci_after
auto_ci_after: clean fmt protobuf
	@git add *.go *.sh *.proto
	@git status
	@files=$$(git status -suno);if [ -n "$$files" ]; then \
		  git status; \
		  git commit -a -m "auto ci [ci-skip]"; \
		  git push origin HEAD:$(branch); \
		  fi;

.PHONY: auto_ci
auto_fmt := find . -name '*.go' -not -path './vendor/*' | xargs goimports -l -w
auto_ci: clean fmt_proto fmt_shell protobuf
	@-find . -name '*.go' -not -path './vendor/*' | xargs gofmt -l -w -s
	@-${auto_fmt}
	@-find . -name '*.go' -not -path './vendor/*' | xargs gofmt -l -w -s
	@${auto_fmt}
	@git add -u
	@git status
	@files=$$(git status -suno); if [ -n "$$files" ]; then \
		  git add -u; \
		  git status; \
		  git commit -a -m "auto ci"; \
		  git remote add originx $(originx); \
		  git remote -v; \
		  git push --quiet --set-upstream originx HEAD:$(branch); \
		  git log -n 2; \
		  exit 1; \
		  fi;


addupstream:
	git remote add upstream git@github.com:assetcloud/plugin.git
	git remote -v

sync:
	git fetch upstream
	git checkout master
	git merge upstream/master
	git push origin master

branch:
	make sync
	git checkout -b ${b}

push:
	@if [ -n "$$m" ]; then \
	git commit -a -m "${m}" ; \
	fi;
	make sync
	git checkout ${b}
	git merge master
	git push origin ${b}

pull:
	@remotelist=$$(git remote | grep ${name});if [ -z $$remotelist ]; then \
		echo ${remotelist}; \
		git remote add ${name} git@github.com:${name}/plugin.git ; \
	fi;
	git fetch ${name}
	git checkout ${name}/${b}
	git checkout -b ${name}-${b}
pullsync:
	git fetch ${name}
	git checkout ${name}-${b}
	git merge ${name}/${b}
pullpush:
	@if [ -n "$$m" ]; then \
	git commit -a -m "${m}" ; \
	fi;
	make pullsync
	git push ${name} ${name}-${b}:${b}

webhook_auto_ci: clean fmt_proto fmt_shell protobuf
	@-find . -name '*.go' -not -path './vendor/*' | xargs gofmt -l -w -s
	@-${auto_fmt}
	@-find . -name '*.go' -not -path './vendor/*' | xargs gofmt -l -w -s
	@${auto_fmt}
	@git status
	@files=$$(git status -suno);if [ -n "$$files" ]; then \
		  git status; \
		  git commit -a -m "auto ci"; \
		  git push origin ${b}; \
		  exit 0; \
		  fi;

webhook:
	git checkout ${b}
	make webhook_auto_ci name=${name} b=${b}
