build_dir := out
config_file := api.toml
BINARY := content-api

DEV_OUT := $(build_dir)/$(BINARY)
LINUX_OUT := $(build_dir)/linux/$(BINARY)

LOCAL_CONFIG_FILE := $(HOME)/config/$(config_file)

VERSION := `git describe --tags`
BUILD := `date +%FT%T%z`
COMMIT := `git log --max-count=1 --pretty=format:%aI_%h`

LDFLAGS := -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD} -X main.commit=${COMMIT}"

BUILD_LINUX := GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(LINUX_OUT) -v .

.PHONY: build linux deploy lastcommit mkbuild clean
# Development
dev :
	go build $(LDFLAGS) -o $(DEV_OUT) -v .

# Run development build
run :
	./$(DEV_OUT)

# Cross compiling linux on for dev.
linux :
	$(BUILD_LINUX)

# From local machine to production server
# Copy env varaible to server
config :
	rsync -v $(LOCAL_CONFIG_FILE) tk11:/home/node/config

deploy : config linux
	rsync -v $(LINUX_OUT) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart $(BINARY)

# For CI/CD
build :
	gvm install go1.13.4
	gvm use go1.13.4
	$(BUILD_LINUX)

downconfig :
	rsync -v tk11:/home/node/config/$(config_file) $(HOME)/config

publish :
	rsync -v $(LINUX_OUT) $(HOME)/go/bin

restart :
	sudo supervisorctl restart $(BINARY)

clean :
	go clean -x
	rm -rf $(build_dir)/*
