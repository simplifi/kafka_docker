default : test

test : build
	go test -v ./...

build :
	go mod tidy
	go build

lint :
	golint ./...
