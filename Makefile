SRC = $(wildcard *.go)
SRC += $(wildcard */*.go)


.PHONY: all run clean

all: strlang

run: strlang
	rlwrap ./strlang

strlang: $(SRC)
	go build .

clean:
	rm -f ./strlang
