version := $(shell git describe --tags)
commit := $(shell git rev-parse --short HEAD)
commit_long := $(shell git rev-parse HEAD)
branch := $(shell git rev-parse --abbrev-ref HEAD)
curtime := $(shell date "+%Y-%m-%d %T %Z")

BINARY=mirage
IMPORT_PATH=github.com/OdyseeTeam/mirage
LDFLAGS="-s -w -X ${IMPORT_PATH}/internal/version.version=$(version) -X ${IMPORT_PATH}/internal/version.commit=$(commit) -X ${IMPORT_PATH}/internal/version.commitLong=$(commit_long) -X ${IMPORT_PATH}/internal/version.branch=$(branch) -X '${IMPORT_PATH}/internal/version.date=$(curtime)'"

.PHONY: test
test:
	go test -cover ./...

linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o dist/linux_amd64/${BINARY} -ldflags ${LDFLAGS}

macos:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o dist/darwin_amd64/${BINARY} -ldflags ${LDFLAGS}

version := $(shell git describe --abbrev=0 --tags|sed 's/v//')
cur_branch := $(shell git rev-parse --abbrev-ref HEAD)
.PHONY: image
image:
	docker buildx build -t odyseeteam/${BINARY}:$(version) -t odyseeteam/${BINARY}:latest -t odyseeteam/${BINARY}:$(cur_branch) --platform linux/amd64 .

version := $(shell git describe --abbrev=0 --tags|sed 's/v//')
.PHONY: publish_image
publish_image:
	docker push odyseeteam/${BINARY}:$(version)
	docker tag odyseeteam/${BINARY}:$(version) odyseeteam/${BINARY}:latest odyseeteam/${BINARY}:$(cur_branch)
	docker push odyseeteam/${BINARY}:latest

tag := $(shell git describe --abbrev=0 --tags)
.PHONY: retag
retag:
	@echo "Re-setting tag $(tag) to the current commit"
	git tag -d $(tag)
	git tag $(tag)