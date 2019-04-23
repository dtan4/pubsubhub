NAME := pubsubhub

SRCS := $(shell find . -type f -name '*.go')

bin/$(NAME): $(SRCS)
	GO111MODULE=on go build -o bin/$(NAME) github.com/dtan4/pubsubhub/cmd/pubsubhub
