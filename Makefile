BINARY_NAME=supportbot

init:
	go mod init ${BINARY_NAME}

run:
	clear
	go test -timeout 30s -v -cover
	go run .

test:
	clear
	go test -timeout 30s -v -cover

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags "-s -w" -o build/${BINARY_NAME}-darwin_x64
	GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o build/${BINARY_NAME}-linux_x64
	GOARCH=amd64 GOOS=windows go build -ldflags "-s -w" -o build/${BINARY_NAME}-windows.exe
	GOARCH=arm GOOS=linux GOARM=5 go build -ldflags "-s -w" -o build/${BINARY_NAME}-linux_arm_5
	GOARCH=arm GOOS=linux GOARM=7 go build -ldflags "-s -w" -o build/${BINARY_NAME}-linux_arm_7
	GOARCH=arm64 GOOS=linux GOARM=7 go build -ldflags "-s -w" -o build/${BINARY_NAME}-linux_arm64_7

clean:
	rm -fr build

cacheclean:
	go clean -cache
	go clean -testcache
	go clean -modcache

dep:
	go mod download

vet:
	go vet

all: clean build