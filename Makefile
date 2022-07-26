SHELL=/usr/bin/env bash

all: build
.PHONY: all

BINS:=

deal-sync:
	rm -f deal-sync
	go build -mod=vendor -o deal-sync main.go
BINS+=deal-sync


# ieth-cmd:
# 	rm -rf ieth-cmd
# 	go build -mod=vendor -o ieth-cmd cmd/main.go
# BINS+=ieth-cmd

build: $(BINS)


clean:
	rm -f $(BINS)
