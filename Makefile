GHACCOUNT := hooklift
NAME := nnn
VERSION := v0.1.0

OK_COLOR	:= \x1b[32;01m
NO_COLOR	:= \x1b[0m

include xhyve.mk

build: libxhyve.a
	@echo "$(OK_COLOR)------> Running go build...$(NO_COLOR)"
	go build

# This is so the linker doesn't complain given that CGO is already defining a
# main function.
patch-xhyve:
	cd vendor/xhyve/src; patch < ../../../xhyve.c.patch

.PHONY: build
