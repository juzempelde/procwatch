all: procwatch procwatch-win

procwatch:
	cd backend; go build -v -o ../procwatch ./cmd/procwatch

procwatch-win:
	cd backend; GOOS=windows GOARCH=386 go build -v -o ../procwatch_x86.exe ./cmd/procwatch
	cd backend; GOOS=windows GOARCH=amd64 go build -v -o ../procwatch_amd64.exe ./cmd/procwatch

test:
	cd backend; go test ./...

test-cover:
	cd backend; go test -cover ./...

.PHONY: \
	procwatch \
	procwatch-win \
	test
