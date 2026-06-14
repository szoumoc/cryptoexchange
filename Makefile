build:
	go build -o bin/ ./...
run:
	go run ./...
test:
	go test -v ./...

test-git:
	go test -v ./... && make gitorderbook

gitorderbook:
	git add .
	git commit -m "update orderbook and tests"
	git push