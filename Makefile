
SRC = $(wildcard */*.go)

all: strlang

run: strlang
	rlwrap ./strlang

strlang: $(SRC)
	go build .
