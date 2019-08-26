test : build
	go test ./...
build : get
	go build
get :
	go mod download

