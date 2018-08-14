glide_install:
	@echo "Vendoring dependencies"
	@glide install

fixfmt:
	@echo "Fixing format"
	@gofmt -w -l .

test:
	@echo "Running Tests"
	@go test ./...

compile:
	@echo "Building Binaries"
	@GOOS=darwin GOARCH=amd64 go build -o=build/godownload-darwin app/main.go
	@GOOS=linux GOARCH=amd64 go build -o=build/godownload-linux app/main.go
	@GOOS=windows GOARCH=amd64 go build -o=build/godownload-windows app/main.go

build:	clean	compile	test

clean:
	@echo "Removing existing builds"
	@-rm -rf build

checkfmt:
	@bash -c "diff -u <(echo -n) <(gofmt -d ./)"

