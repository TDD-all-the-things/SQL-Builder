.PHONY: test
test:
	@go test -race ./...

.PHONY: testv
testv:
	@go test -v -count=1 -race ./...

.PHONY: setup
setup:
	@sh ./script/setup.sh

.PHONY: lint
lint:
	@golangci-lint run

.PHONY:	fmt
fmt:
	@goimports -l -w .

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) fmt
	@$(MAKE) tidy
	@$(MAKE) lint