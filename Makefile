help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

depupdate:
	dep ensure -update

build:
	go build -o terraform-provider-searchstax

install: build # Install searchstax provider into terraform plugin directory
	mv terraform-provider-searchstax ~/.terraform.d/plugins/

init: install # Run terraform for local testing
	terraform init

.PHONY: help build install init
.DEFAULT_GOAL := help