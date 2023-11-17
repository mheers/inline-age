all: build

build: build-mac build-linux

build-mac:
	GOOS=darwin GOARCH=arm64 go build -o ia-darwin-arm64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ia

prepare-tests:
	docker-compose -f reference/git/test/docker-compose.yml pull
	docker-compose -f reference/vault/test/docker-compose.yml pull

git-start: ## Starts a local git server
	docker-compose -f reference/git/test/docker-compose.yml -f reference/git/test/docker-compose.dev.yml up

git-init: ## Initializes the local git server
	docker-compose -f reference/git/test/docker-compose.yml -f reference/git/test/docker-compose.dev.yml exec git sh /bin/init.sh


vault-start: ## Starts a local vault server
	docker-compose -f reference/vault/test/docker-compose.yml -f reference/vault/test/docker-compose.dev.yml up

vault-init: ## Initializes the local vault server
	docker-compose -f reference/vault/test/docker-compose.yml -f reference/vault/test/docker-compose.dev.yml exec vault sh /bin/init.sh
