.PHONY: all
all: clean build test
TAGS = awskms awssecretsmanager awsssm

clean:
	rm ./secure-exec-*amd64 || true

test:
	go test -v -tags "$(TAGS)" ./... -coverprofile=coverage.txt -covermode=atomic

build:
	GOOS=linux GOARCH=amd64 go build -i -tags '$(TAGS)' -o "secure-exec-linux-amd64"
	GOOS=darwin GOARCH=amd64 go build -i -tags '$(TAGS)' -o "secure-exec-darwin-amd64"

docker:
	docker build -t secure-exec-example .
