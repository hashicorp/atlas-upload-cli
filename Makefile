NAME = $(shell cat ./main.go | grep "Name = " | cut -d" " -f4 | sed 's/[^"]*"\([^"]*\).*/\1/')
VERSION = $(shell cat ./main.go | grep "Version = " | cut -d" " -f4 | sed 's/[^"]*"\([^"]*\).*/\1/')
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: build

deps: Makefile
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d
	touch $@

build: bin/$(NAME)

bin/$(NAME): deps
	@mkdir -p bin/
	go build -o bin/$(NAME)

test: deps
	go list ./... | xargs -n1 go test -timeout=3s

xcompile: deps test
	@rm -rf build/
	@mkdir -p build
	gox \
		-output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -zcvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done

clean:
	rm -rf bin build deps

.PHONY: all build test xcompile package clean
