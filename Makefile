SERVER_ID=server1
PROFILE=larwef
REGION=eu-west-1
ACCOUNT_ID=$$(aws sts get-caller-identity --output text --query 'Account' --profil $(PROFILE))
ECS_REPO=$(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/grpc-test

# Collections of commands
build: build-linux build-mac build-windows build-docker
docker: build-docker run-docker

# Run locally
run-server:
	serverId=$(SERVER_ID) go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

run-docker:
	docker run -it --rm -p 8080:8080 -e serverId=$(SERVER_ID) go-grpc-test-server

# Generate grpc code
proto:
	protoc -I internal/hello/ internal/hello/*.proto --go_out=plugins=grpc:internal/hello

# Build
build-linux:
	GOOS=linux go build -o target/linux/server cmd/server/main.go
	GOOS=linux go build -o target/linux/client cmd/client/main.go

build-mac:
	GOOS=darwin go build -o target/mac/server cmd/server/main.go
	GOOS=darwin go build -o target/mac/client cmd/client/main.go

build-windows:
	GOOS=windows go build -o target/windows/server cmd/server/main.go
	GOOS=windows go build -o target/windows/client cmd/client/main.go

build-docker:
	docker build -t go-grpc-test-server -f build/docker/Dockerfile .

# Upload to ECS repository
upload-docker:
	docker tag go-grpc-test-server $(ECS_REPO)
	$$(aws ecr get-login --no-include-email --profile $(PROFILE) --region $(REGION))
	docker push $(ECS_REPO)