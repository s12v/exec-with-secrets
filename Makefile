.PHONY: all
all: clean test build
TAGS = awskms awssecretsmanager awsssm

clean:
	rm -rf ./bin || true

test:
	go test -v -tags "$(TAGS)" ./... -coverprofile=coverage.txt -covermode=atomic

build:
	GOOS=linux GOARCH=amd64 go build -i -tags '$(TAGS)' -ldflags='-s -w' -o "bin/secure-exec-linux-amd64"
	GOOS=darwin GOARCH=amd64 go build -i -tags '$(TAGS)' -ldflags='-s -w' -o "bin/secure-exec-darwin-amd64"

docker:
	docker build -t secure-exec-example .
