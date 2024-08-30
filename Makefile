test:
	go test internal/app/handler/*

serve:
	go run cmd/shortener/main.go -a localhost:8081 -b http://localhost:8081 -d "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_sh sslmode=disable"

build:
	go build -o shortener *.go