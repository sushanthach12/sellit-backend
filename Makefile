# '@' is used to suppress the command output in the terminal.


# .PHONY is used to declare phony targets, which are not actual files but just names for commands to be executed.
# This is useful to avoid re-build failure, since make will consider a target up-to-date if a file with the same name exists in the directory.
# 'make: build is up to date' error occurs when a file with the same name as the target exists in the directory, 
# and make assumes that the target is already built and up-to-date.
.PHONY: build run migrate-up migrate-down

build:
	@go build -o bin/api ./cmd/api

run: build
	@./bin/api

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down
