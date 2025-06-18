LOG_LEVEL ?= disabled

run:
	go run ./main.go -log-level $(LOG_LEVEL)