config_file_name := api.toml
local_config_file := $(HOME)/config/api.toml

app_name := content-api
go_version := go1.19

current_dir := $(shell pwd)
sys := $(shell uname -s)
hardware := $(shell uname -m)
src_dir := $(current_dir)
build_dir := $(current_dir)/build

default_exec := $(build_dir)/$(sys)/$(hardware)/$(app_name)

linux_x86_exec := $(build_dir)/linux/x86/$(app_name)

linux_arm_exec := $(build_dir)/linux/arm/$(app_name)

.PHONY: build
build : version
	go build -o $(default_exec) -tags production -v $(src_dir)

.PHONY: builddir
builddir :
	mkdir -p $(build_dir)

# Create version info to be embedded in go binary.
.PHONY: version
version : builddir
	git describe --tags > build/version
	git log --max-count=1 --pretty=format:%aI_%h > build/commit
	date +%FT%T%z > build/build_time

.PHONY: devenv
devenv : builddir
	rsync $(HOME)/config/env.dev.toml $(build_dir)/$(config_file_name)

# Run in dev mode
.PHONY: run
run :
	$(default_exec)

# Run in production mode.
.PHONY: prod
prod :
	./$(build_dir)/$(BINARY) -production

.PHONY: install-go
install-go:
	@echo "* Install go version $(go_version)"
	gvm install $(go_version)

# Sync env file to build directory.
.PHONY: config
config : builddir
	@echo "* Pulling config  file from server"
	# Download configuration file
	rsync -v node@tk11:/home/node/config/$(config_file_name) $(build_dir)/$(config_file_name)

# Build linux x86 binary
# In Jenkins, ensure to run `config` before this one.
.PHONY: amd64
amd64 : version
	@echo "Build production linux version $(version)"
	GOOS=linux GOARCH=amd64 go build -o $(linux_x86_exec) -tags production -v $(src_dir)

# Build linux arm binary
.PHONY: arm
arm : version
	@echo "Build production arm version $(version)"
	GOOS=linux GOARM=7 GOARCH=arm go build -o $(linux_arm_exec) -tags production -v $(src_dir)

# Sync binary to production server
.PHONY: publish
publish :
	rsync -v $(linux_x86_exec) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart $(BINARY)

.PHONY: clean
clean :
	go clean -x
	rm build/*


