deps:
	go get -d -t ./...
test:
	go test -v ./...
install: deps
	go install

build:
	gox -output=pkg/dist/mysqldiff_{{.OS}}_{{.Arch}}
