.PHONY: docker-up docker-down run-http run-worker run-cron run-go

docker-up:
	docker-compose -f docker-compose-development.yml up -d

docker-down:
	docker-compose -f docker-compose-development.yml down

run-http:
	go run cmd/http/main.go

run-worker:
	go run cmd/worker/main.go

run-cron:
	go run cmd/cron/main.go

test:
	go test ./...