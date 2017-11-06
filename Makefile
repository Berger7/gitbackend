default: test

pre_env:
	./prepare_env.sh

build: pre_env
	go build

test: pre_env
	go test -race -coverprofile=coverage.txt -covermode=atomic