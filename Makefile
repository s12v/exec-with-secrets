.PHONY: all
all: deps clean test build

deps:
	dep ensure -v

clean: 
	rm ./secure-exec || true

test:
	go test -v -tags "awskms awssecretsmanager awsssm" ./... -coverprofile=coverage.txt -covermode=atomic

build:
	GOOS=linux GOARCH=amd64 go build -i -tags 'awskms awssecretsmanager awsssm' -o "secure-exec-linux-amd64"
	GOOS=darwin GOARCH=amd64 go build -i -tags 'awskms awssecretsmanager awsssm' -o "secure-exec-darwin-amd64"

docker:
	docker build -t secure-exec-example .
