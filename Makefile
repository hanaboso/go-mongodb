DC=docker-compose
DE=docker-compose exec -T app
IMAGE=dkr.hanaboso.net/hanaboso/go-mongodb

.env:
	sed -e 's/{DEV_UID}/$(shell id -u)/g' \
		-e 's/{DEV_GID}/$(shell id -g)/g' \
		-e 's|{GITLAB_CI}|$(shell [ ! -z "$$GITLAB_CI" ] && echo true || echo false)|g' \
		-e 's|{DOCKER_SOCKET_PATH}|$(shell test -S /var/run/docker-$${USER}.sock && echo /var/run/docker-$${USER}.sock || echo /var/run/docker.sock)|g' \
		.env.dist >> .env; \

docker-up-force: .env
	$(DC) pull
	$(DC) up -d --force-recreate --remove-orphans

docker-down-clean: .env
	$(DC) down -v

docker-compose.ci.yml:
	# Comment out any port forwarding
	sed -r 's/^(\s+ports:)$$/#\1/g; s/^(\s+- \$$\{DEV_IP\}.*)$$/#\1/g; s/^(\s+- \$$\{GOPATH\}.*)$$/#\1/g' docker-compose.yml > docker-compose.ci.yml

go-update:
	$(DE) su-exec root go get -u all
	$(DE) su-exec root go mod tidy
	$(DE) su-exec root chown dev:dev go.mod go.sum

init-dev: docker-up-force
	$(DE) go mod download

lint:
	$(DE) gofmt -w .
	$(DE) golint ./...

fast-test: lint
	$(DE) mkdir var || true
	$(DE) go test -cover -coverprofile var/coverage.out ./... -count=1
	$(DE) go tool cover -html=var/coverage.out -o var/coverage.html

test: init-dev fast-test docker-down-clean
