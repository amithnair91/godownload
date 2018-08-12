glide_install:
			glide install

fixfmt:
			gofmt -w -l .

test:
			 go test ./...

