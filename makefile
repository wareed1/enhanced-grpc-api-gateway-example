BUF_VERSION:=v1.9.0
SWAGGER_UI_VERSION:=v4.15.5

clean:
	rm -rf ./out
	rm -rf ./third_party/OpenAPI
	rm -rf ./gen

lint:
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) lint

generate: generate/proto generate/swagger-ui
generate/proto:
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) generate
generate/swagger-ui:
	SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) /bin/bash ./scripts/generate-swagger-ui.sh

secure-tls:
ifeq ($(origin IPADDR),undefined)
	@echo "IPADDR not set"
	exit 1
endif
	mkdir -p certs && (cd certs; /bin/bash ../scripts/mk-certs.sh; cd ..)
	mkdir -p ./secure && /bin/bash ./scripts/mk-secure.sh
	@echo "modifying source" && /bin/bash ./scripts/change-source.sh secure

insecure-tls:
	@echo "returning original source" && /bin/bash ./scripts/change-source.sh insecure

build:
	mkdir -p ./out
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/apigw ./api-gw
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/users ./users
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/orders ./orders
	docker-compose build

run-servers:
	@echo "--> Starting servers"
	@docker-compose up
