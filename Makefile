GHACCOUNT := hooklift
NAME := nnn
VERSION := v0.1.0

OK_COLOR	:= \x1b[32;01m
NO_COLOR	:= \x1b[0m

include xhyve.mk

CGO_LDFLAGS 	:= -ldflags "-X main.Version=$(VERSION) -X main.Name=$(NAME)"

build: libxhyve.a
	@echo "$(OK_COLOR)------> Running go build...$(NO_COLOR)"
	go build -x $(CGO_LDFLAGS) -o $(NAME)

# This is so the linker doesn't complain given that CGO is already defining a
# main function.
patch-xhyve:
	cd vendor/xhyve; patch -p1 < ../../xhyve.patch

.PHONY: build
