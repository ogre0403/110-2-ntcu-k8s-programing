IMG ?= ogre0403/incluster
TAG ?= latest


run-client:
	go run ./cmd/client

run-in-cluster:
	go run ./cmd/incluster -outside-cluster

run-informer:
	go run ./cmd/informer -outside-cluster


#build-in-cluster-image:
#	docker build -t $(IMG):$(TAG) .
#
#load-in-cluster-image: build-in-cluster-image
#	kind load docker-image $(IMG):$(TAG)


deploy-informer:
	docker build --build-arg SAMPLE=informer  -t $(IMG):$(TAG) .
	kind load docker-image $(IMG):$(TAG)
	kubectl apply -f  manifest/deployment.yaml -f manifest/rbac.yaml

deploy-incluster:
	docker build --build-arg SAMPLE=incluster  -t $(IMG):$(TAG) .
	kind load docker-image $(IMG):$(TAG)
	kubectl apply -f  manifest/deployment.yaml -f manifest/rbac.yaml