test:
	@ go test ./...

auto_test:
	@ fswatch -or . | xargs -n1 -I{} go test ./...

run:
	@ go run main.go

format:
	@ go fmt ./...

build:
	@ docker compose build

up:
	@ docker compose up -d

down:
	@ docker compose down

connect:
	@ docker compose exec goapp bash

# release:
# 	@
