#!/usr/bin/make -f

GO=go
GB=gb

darwin:
	env GOOS=darwin GOARCH=amd64 $(GB) build

linux:
	env GOOS=linux GOARCH=amd64 $(GB) build

all: test clean darwin linux

clean:
	rm -fR pkg bin

test:
	$(GB) test -test.v
