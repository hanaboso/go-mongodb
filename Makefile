DC=docker-compose
IMAGE=dkr.hanaboso.net/pipes/pipes/mongodb

.env:
	sed -e 's/{DEV_UID}/$(shell id -u)/g' \
		-e 's/{DEV_GID}/$(shell id -u)/g' \
		-e 's/{SSH_AUTH}/$(shell if [ '$(shell uname)' = 'Linux' ]; then echo '\/tmp\/.ssh-auth-sock'; else echo '\/tmp\/.nope'; fi)/g' \
		.env.dist >> .env; \

build:
	docker build -t ${IMAGE}:${TAG} .
	docker push ${IMAGE}:${TAG}

docker-up-force: .env
	$(DC) pull
	$(DC) up -d --force-recreate --remove-orphans

docker-down-clean: .env
	$(DC) down -v

init-dev: docker-up-force
