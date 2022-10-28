CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/cmd/bot
DB_USER=${FINANCIAL_BOT_DB_USER}
DB_PASS=${FINANCIAL_BOT_DB_PASS}
DB_HOST=${FINANCIAL_BOT_DB_HOST}
DB_PORT=${FINANCIAL_BOT_DB_PORT}
DB_NAME=${FINANCIAL_BOT_DB_NAME}

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

generate: install-mockgen
	${MOCKGEN} -source=internal/model/messages/incoming_msg.go -destination=internal/mocks/messages/incoming_msg.go
	${MOCKGEN} -source=internal/model/callbacks/incoming_callback.go -destination=internal/mocks/callbacks/incoming_callback.go
	${MOCKGEN} -source=internal/service/calculator_service.go -destination=internal/mocks/service/calculator_service.go
	${MOCKGEN} -source=internal/service/currency_exchange_service.go -destination=internal/mocks/service/currency_exchange_service.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

docker-run:
	sudo docker compose up

migrate:
	goose -dir migrations postgres "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
