all: build

build: arachne ivd kubernetes rest_api

cmd:
	cd cmd/server; go get; go install

arachne:
	cd pkg/arachne; go get; go install

ivd:
	cd pkg/ivd; go get; go install

kubernetes:
	cd pkg/kubernetes; go get; go install

rest_api:
	cd pkg/rest_api; go get; go install
