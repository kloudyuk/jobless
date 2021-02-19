SHELL = /bin/bash

PROFILE = kloudy
REGION = eu-west-2
IMAGE = public.ecr.aws/l7t2e6a6/jobless

.PHONY: image
image:
	go mod tidy
	docker build -t $(IMAGE):latest .

.PHONY: push
push: image
	docker push $(IMAGE):latest
