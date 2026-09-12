.PHONY: all clean test build docker
all: clean test build
TAGS = awskms awssecretsmanager awsssm
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

clean:
	rm -rf ./bin || true

test:
	go test -v -tags "$(TAGS)" ./... -coverprofile=coverage.txt -covermode=atomic

build:
	@mkdir -p bin
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		echo "building exec-with-secrets-$$os-$$arch"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -tags '$(TAGS)' -ldflags='-s -w' -o "bin/exec-with-secrets-$$os-$$arch"; \
	done

docker:
	docker build -t exec-with-secrets-example .
