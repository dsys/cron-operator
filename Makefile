.PHONY: all
all: build

.PHONY: clean
clean:
	@ rm -rf build
	@ mkdir -p build

.PHONY: publish
build: TAG ?= latest
build: clean
	@ go build cmd/operator/main.go
	@ ./node_modules/.bin/webpack -p --output-path build/www
	@ docker build -t pavlov/cron-operator:$(TAG) .

.PHONY: push
push: TAG ?= latest
push:
	@ docker push pavlov/cron-operator:$(TAG)

.PHONY: build-webpack
build-webpack:
	@ ./node_modules/.bin/webpack --output-path build/www

.PHONY: watch-webpack
watch-webpack:
	@ ./node_modules/.bin/webpack --output-path build/www --watch

.PHONY: proxy
proxy: build-webpack
	@ kubectl proxy --www build/www --www-prefix=/ --api-prefix=/k8s-api
