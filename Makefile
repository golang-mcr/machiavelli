NAME = machiavelli

all: build-sender build-receiver
	
build-sender:
	go build -o $(NAME)-sender -tags sender

build-receiver:
	go build -o $(NAME)-receiver -tags receiver

listen:
	./machiavelli-receiver --config config.gcfg
send:
	./machiavelli-sender --config config.gcfg --message "Hello world!"

cs-fix:
	go fmt .
