all: procwatch

procwatch:
	go build -o procwatch ./backend/cmd/procwatch

test:
	go test ./backend/...

test-cover:
	go test -cover ./backend/...

.PHONY: \
	procwatch \
	test
