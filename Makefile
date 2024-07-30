test:
	go test internal/app/handler/*

serve:
	go run cmd/shortener/main.go -a localhost:8081 -b http://localhost:8081

build:
	go build -o shortener *.go