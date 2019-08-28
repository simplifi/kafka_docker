default : test

test : build
	go test ./...

build :
	go mod tidy
	go build

lint :
	golint ./...
