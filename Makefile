NAME = machiavelli

setup:
	cp config.gcfg.example config.gcfg

all: build-sender build-receiver
	
build-sender:
	go build -o $(NAME)-sender -tags sender

build-receiver:
	go build -o $(NAME)-receiver -tags receiver
