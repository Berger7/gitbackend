default: test

pre_env:
	./prepare_env.sh

install: pre_env
	glide install

build: install
	go build

test: install
	go test