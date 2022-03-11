.PHONY: build
build:
	go build -o myhttp -race main.go

.PHONY: test
test:
	go test -v -race -coverprofile cover.out; \
	go tool cover -func cover.out

.PHONY: example
example:
	make build; \
	./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
