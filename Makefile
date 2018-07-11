all: procwatch procwatch-win

procwatch:
	go build -o procwatch ./backend/cmd/procwatch

procwatch-win:
	GOOS=windows GOARCH=386 go build -v -o procwatch_x86.exe ./backend/cmd/procwatch
	GOOS=windows GOARCH=amd64 go build -v -o procwatch_amd64.exe ./backend/cmd/procwatch

test:
	go test ./backend/...

test-cover:
	go test -cover ./backend/...

.PHONY: \
	procwatch \
	procwatch-win \
	test
