PROJECTNAME=$(shell basename "$(PWD)")
BINLOC =./dist/service-apigw

## help: List available command
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## generate.sh-pb: Generate gRPC stubs.
generate-pb:
	@sh generate.sh service_auth
	@sh generate.sh service_area

## build-cmd: Build CMD.
build-cmd:
	go build -o $(BINLOC) -gcflags="all=-N -l" -ldflags="-X github.com/luraproject/lura/v2/core.KrakendVersion=MAKIMA" ./cmd

## build-plugins: Build all plugins as Go Plugins (.so).
build-plugins: cue-export
	mkdir -p dist/plugins
	./auto-generator.sh

build-ci: generate.sh build-cmd build-plugins

# run: Run gateway
run:
	go run ./cmd

## start: Start program
start: build-cmd build-plugins run

## cue-export: export CUE configuration files as JSON
cue-export: cue-export-prod
	cue export ./config/endpoint -f -o cmd/generated-lura-conf.json

## cue-export-prod: export CUE configuration files as JSON (for NTT production deployment)
cue-export-prod:
	cue export ./config/endpoint  -f -o cmd/generated-lura-conf.json

## cue-eval: validate the correctness of CUE configuration files
cue-eval:
	cue eval ./config/endpoint

## cue-format: idiomatically format CUE configuration files
cue-format:
	cue fmt ./config -E -v --strict
	cue fmt ./config/eva -E -v --strict
	cue fmt ./config/endpoint -E -v --strict
