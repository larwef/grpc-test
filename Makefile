APP_NAME=go-grpc-test-server
VERSION=0.0.4
GOOS=linux
PORT=8080
PROFILE=larwef
REGION=eu-west-1
ACCOUNT_ID=$$(aws sts get-caller-identity --output text --query 'Account' --profil $(PROFILE))
ECS_REPO=$(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/$(APP_NAME)

# Collections of commands
build: build-server build-docker
docker: build-docker run-docker

# Run locally
run-server:
	port=$(PORT) go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

run-docker:
	docker run -it --rm -p $(PORT):$(PORT) \
	-e port=$(PORT) \
	$(APP_NAME)

# Generate grpc code
proto:
	protoc -I internal/hello/ internal/hello/*.proto --go_out=plugins=grpc:internal/hello

# Build
build-server:
	GOOS=$(GOOS) go build -ldflags "-X main.version=$(VERSION)" -o target/server/app cmd/server/main.go

build-docker:
	docker build -t $(APP_NAME) -f build/docker/Dockerfile .

# Upload to ECS repository
upload-docker:
	docker tag $(APP_NAME) $(ECS_REPO):$(VERSION)
	docker tag $(APP_NAME) $(ECS_REPO):latest
	$$(aws ecr get-login --no-include-email --profile $(PROFILE) --region $(REGION))
	docker push $(ECS_REPO)

# Test caching is useful, but dont want it for these tests. Using non-cacheable flag.
# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go test -v -count=1 ./...