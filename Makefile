.PHONY: fmt compile run 

PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
GOBASE=$(shell pwd)

GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will store the server process id when it's running on development mode
PID=/tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

fmt:
	gofmt -w .

compile:
	@for GOOS in linux windows;do\
		for GOARCH in amd64; do\
			GOOS=$${GOOS} GOARCH=$${GOARCH} go build -v -o  bin/$(PROJECTNAME).$${GOOS}-$${GOARCH} ; \
		done ;\
	done
	cp bin/go-http-frame.windows-amd64 bin/go-http-frame
	-rm bin/go-http-frame.windows-amd64
	-rm bin/go-http-frame.linux-amd64
	
run:
	go run main.go
