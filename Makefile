run-server:
	go run cmd/server/main.go 'server1'

run-client:
	go run cmd/client/main.go

proto:
	protoc -I internal/hello/ internal/hello/*.proto --go_out=plugins=grpc:internal/hello

build: build-linux build-mac build-windows

build-linux:
	GOOS=linux go build -o target/linux/server cmd/server/main.go
	GOOS=linux go build -o target/linux/client cmd/client/main.go

build-mac:
	GOOS=darwin go build -o target/mac/server cmd/server/main.go
	GOOS=darwin go build -o target/mac/client cmd/client/main.go

build-windows:
	GOOS=windows go build -o target/windows/server cmd/server/main.go
	GOOS=windows go build -o target/windows/client cmd/client/main.go