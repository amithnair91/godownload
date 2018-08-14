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
	@go build -o=build/go-downloader app/main.go

build:	clean	compile	test

clean:
	@echo "Removing existing builds"
	@-rm -rf build

checkfmt:
	@bash -c "diff -u <(echo -n) <(gofmt -d ./)"

