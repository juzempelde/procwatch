all: procwatch

procwatch:
	go build ./backend/cmd/procwatch

test:
	go test ./backend/...

test-cover:
	go test -cover ./backend/...

.PHONY: \
	procwatch \
	test
