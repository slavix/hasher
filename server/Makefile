
run:
	@go run ./cmd/app/main.go

generate:
	@mkdir -p ./internal/generated
	@swagger generate server -f ./api/api.yml -t ./internal/generated --exclude-main

