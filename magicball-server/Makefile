APP_NAME=magicball-server
REGISTRY=docker.io
NAMESPACE=gbsal

build:
	docker build -t $(APP_NAME) .

tag:
	docker tag $(APP_NAME) $(REGISTRY)/$(NAMESPACE)/$(APP_NAME):latest

push:
	docker push $(REGISTRY)/$(NAMESPACE)/$(APP_NAME):latest
