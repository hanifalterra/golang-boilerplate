.PHONY: run-http run-worker run-cron

run-http:
	go run cmd/http/main.go

run-worker:
	go run cmd/worker/main.go

run-cron:
	go run cmd/cron/main.go