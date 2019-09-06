build_dir := out
BINARY := content-api

VERSION := `git describe --tags`
BUILD := `date +%FT%T%z`

LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}"

.PHONY: build linux deploy attack report lastcommit mkbuild clean
build :
	go build $(LDFLAGS) -o $(build_dir)/$(BINARY) -v .

run :
	./$(build_dir)/${BINARY}

prod :
	./$(build_dir)/$(BINARY) -production

deploy : linux
	rsync -v $(build_dir)/linux/$(BINARY) node11:/home/node/go/bin/

linux :
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(build_dir)/linux/$(BINARY) -v .

lastcommit :
	git log --max-count=1 --pretty=format:%aI\ %h

mkbuild :
	mkdir -p build

clean :
	go clean -x
	rm build/*
