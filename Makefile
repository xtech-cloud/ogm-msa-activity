APP_NAME := omo-msa-activity
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

.PHONY: build
build: 
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run: 
	./bin/${APP_NAME}

.PHONY: install
install: 
	go install

.PHONY: clean
clean: 
	rm -rf /tmp/msa-activity.db

.PHONY: call
call:
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Subscribe '{"notification":"omo.msa.account.notification"}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Subscribe '{"notification":"omo.msa.application.notification"}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Subscribe '{"notification":"omo.msa.license.notification"}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Fetch '{"count":5}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Fetch '{"offset":1, "count":1}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Unsubscribe '{"notification":"omo.msa.license.notification"}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Channel.Fetch '{"count":5}'
	MICRO_REGISTRY=consul micro call omo.msa.activity Record.Fetch ''

.PHONY: tcall
tcall:
	mkdir -p ./bin
	go build -o ./bin/ ./tester
	./bin/tester

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo-msa-activity:latest
