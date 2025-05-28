.PHONY: seed

seed:
	@go run cmd/seed/seed.go
.PHONY: gen-docs
gen-docs:
	@swag init --parseDependency --parseInternal -g cmd/api/main.go -o ./docs && swag fmt 