NAME=xchenhao/data-transport-service
VERSION=0.1.0
GOPROXY=https://goproxy.cn

build:
	@echo -e 'NAME=${NAME}\nVERSION=${VERSION}\nGOPROXY=${GOPROXY}\n'
	@docker build --force-rm \
		--build-arg goproxy=${GOPROXY} \
		-t ${NAME}:${VERSION} .

version:
	@echo ${VERSION}
