CC=go

test: todotxt_test.go
	$(CC) test

build: todotxt.go
	$(CC) build
