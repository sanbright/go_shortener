test:
	go test internal/app/handler/*

serve:
	go run cmd/shortener/main.go -a localhost:8083 -b http://localhost:8083

build:
	go build -o shortener *.go