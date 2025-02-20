.PHONY: openapi

openapi:
	swag init -g server.go --parseDependency --parseInternal --dir ./internal/server --output ./openapi
	go mod tidy
