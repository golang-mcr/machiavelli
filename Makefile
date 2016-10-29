NAME = machiavelli

build-sender:
	go build -o $(NAME)-sender -tags sender

build-receiver:
	go build -o $(NAME)-receiver -tags receiver
