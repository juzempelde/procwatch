all: procwatch

procwatch:
	go build ./backend/cmd/procwatch

.PHONY: \
	procwatch
