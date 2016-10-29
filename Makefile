NAME = machiavelli

build:
	go build -o $(NAME)

run: build
	./$(NAME)
