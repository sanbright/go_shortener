test:
	go test internal/app/handler/*

serve:
	go run cmd/shortener/main.go -a localhost:8081 -b http://localhost:8081 -d "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=disable"

build:
	go build -o shortener *.go

prof-p:
	 go tool pprof --http=:8082 -seconds=50 http://localhost:8081/debug/pprof/profile

prof-h:
	 go tool pprof --http=:8082 -seconds=50 http://localhost:8081/debug/pprof/heap

coverage:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
gofmt:
	gofmt -w cmd/*
	gofmt -w internal/*

goimports:
	goimports -local "github.com/sanbright/go_shortener" -w cmd/*
	goimports -local "github.com/sanbright/go_shortener" -w internal/*
