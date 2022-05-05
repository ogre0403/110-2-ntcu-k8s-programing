IMG ?= ogre0403/incluster
TAG ?= latest


run-client:
	go run ./cmd/client

run-in-cluster:
	go run ./cmd/incluster -outside-cluster

build-in-cluster-image:
	docker build -t $(IMG):$(TAG) .

load-in-cluster-image:
	kind load docker-image $(IMG):$(TAG)

run-informer:
	go run ./cmd/informer -outside-cluster

