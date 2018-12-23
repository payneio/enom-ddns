repository=quay.io/payneio
container=enomddns
tag=1

.PHONY: run
run: build/enom-ddns
	enom-ddns

build:
	mkdir -p build

build/enom-ddns: main.go build
	go build -o build/enom-ddns main.go

build/linux-amd64/enom-ddns: main.go build/linux-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o build/linux-amd64/enom-ddns main.go

build/linux-amd64:
	mkdir -p build/linux-amd64

build/container: build/linux-amd64/enom-ddns Dockerfile ca-certificates.crt
	docker build -t $(container) .
	mkdir -p build && touch build/container

.PHONY: release
release: build/container
	docker tag $(container) $(repository)/$(container):$(tag)
	docker push $(repository)/$(container):$(tag)

.PHONY: clean
clean:
	rm -rf build
