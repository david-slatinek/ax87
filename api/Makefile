.DEFAULT_GOAL := help

all: testEnv testUtil testDB testServer ## Perform all tests

compile: ## Compile api
	go build -v .

compileClean: ## Compile api with -a flag
	go build -v -a .

testEnv: ## Test api/env
	go test api/env -v --cover

testUtil: ## Test api/util
	go test api/util -v --cover

testDB: ## Test api/db
	go test api/db -v --cover

testServer: ## Test api/server
	go test api/server -v --cover

clean: ## Delete compiled ./api
	rm api

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
