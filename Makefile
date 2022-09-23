config_file_name := api.toml
local_config_file := $(HOME)/config/api.toml

app_name := content-api
go_version := go1.18

current_dir := $(shell pwd)
sys := $(shell uname -s)
hardware := $(shell uname -m)
src_dir := $(current_dir)
out_dir := $(current_dir)/out
build_dir := $(current_dir)/build

default_exec := $(out_dir)/$(sys)/$(hardware)/$(app_name)

linux_x86_exec := $(out_dir)/linux/x86/$(app_name)

linux_arm_exec := $(out_dir)/linux/arm/$(app_name)

.PHONY: build
build : version
	go build -o $(default_exec) -tags production -v $(src_dir)

.PHONY: version
version :
	git describe --tags > build/version
	git log --max-count=1 --pretty=format:%aI_%h > build/commit
	date +%FT%T%z > build/build_time

.PHONY: run
run :
	$(default_exec)

prod :
	./$(build_dir)/$(BINARY) -production

.PHONY: amd64
amd64 :
	@echo "Build production linux version $(version)"
	GOOS=linux GOARCH=amd64 go build -o $(linux_x86_exec) -tags production -v $(src_dir)

.PHONY: arm
arm :
	@echo "Build production arm version $(version)"
	GOOS=linux GOARM=7 GOARCH=arm go build -o $(linux_arm_exec) -tags production -v $(src_dir)

.PHONY: config
config : builddir
	@echo "* Pulling config  file from server"
	# Download configuration file
	rsync -v node@tk11:/home/node/config/$(config_file_name) $(build_dir)/$(config_file_name)

.PHONY: publish
publish :
	rsync -v $(linux_x86_exec) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart $(BINARY)

.PHONY: clean
clean :
	go clean -x
	rm build/*

.PHONY: builddir
builddir :
	mkdir -p $(build_dir)

.PHONY: devenv
devenv : builddir
	rsync $(HOME)/config/env.dev.toml $(build_dir)/$(config_file_name)
