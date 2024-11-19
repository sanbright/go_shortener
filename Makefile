test:
	go test cmd/shortener/*.go
	go test internal/app/handler/*
	go test internal/app/generator/*
	go test internal/config/*

bench:
	go test -bench internal/app/generator/* -benchmem

serve:
	go run cmd/shortener/main.go -a localhost:8081 -b http://localhost:8081 -d "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=disable"

build:
	go build  -ldflags "-X main.buildVersion=1.0.0 -X main.buildDate=$(date +%Y-%m-%d) -X main.buildCommit=$(git rev-parse HEAD)" -o shortener main.go

prof-p:
	 go tool pprof --http=:8082 -seconds=50 http://localhost:8081/debug/pprof/profile

prof-h:
	 go tool pprof --http=:8082 -seconds=50 http://localhost:8081/debug/pprof/heap

coverage:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...

gofmt:
	gofmt -w cmd/*
	gofmt -w internal/*

static:
	staticcheck ./...

goimports:
	goimports -local "github.com/sanbright/go_shortener" -w cmd/*
	goimports -local "github.com/sanbright/go_shortener" -w internal/*

godoc-l:
	sudo cp -r ./ /usr/local/go/src/sanbright/go_shortener

godoc-rm:
	sudo rm -rf /usr/local/go/src/sanbright

godoc:
	sudo cp -r ./ /usr/local/go/src/sanbright/go_shortener
	godoc -http=:9009 --play