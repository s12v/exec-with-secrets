.PHONY: all
all: deps clean test build

deps:
	dep ensure -v

clean: 
	rm ./secure-exec || true

test:
	go test -v -tags "awskms awsssm" ./... -coverprofile=coverage.txt -covermode=atomic

build:
	GOOS=linux GOARCH=amd64 go build -i -tags 'awskms awsssm' -o secure-exec

docker:
	docker build -t secure-exec-example .
