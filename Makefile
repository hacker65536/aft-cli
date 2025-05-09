SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help


.PHNOY: install
install:
	go install

.PHNOY: build
build:
	go build

.PHNOY: clean
clean:
	go clean