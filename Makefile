BINARY_NAME=links-service

build:
	go build -o ${BINARY_NAME}-windows services\links\cmd\main.go

run:
	./${BINARY_NAME}-windows

build_and_run:	build run

test:
	go test ./...

test-with-component:
	go test --tags=component ./...