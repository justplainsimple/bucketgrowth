build:
	go build -o .build/bucketgrowth bucketgrowth/cmd/bucketgrowth

format:
	go fmt ./...

test:
	go test ./...

vet:
	go vet ./...

sast:
	gosec ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
