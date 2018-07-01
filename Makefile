all: procwatch

procwatch:
	go build ./backend/cmd/procwatch

test:
	go test ./backend/...

.PHONY: \
	procwatch \
	test
