glide_install:
			glide install

fixfmt:
			gofmt -w -l .

test:
			 go test ./...

build:
			 go build -o=build/go-downloader app/main.go

