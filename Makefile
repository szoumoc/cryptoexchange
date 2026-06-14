build:
	go build -o bin/ ./...
run:
	go run ./...
test:
	go test -v ./...

git:
	git add .
	git commit -m "update orderbook and tests"
	git push