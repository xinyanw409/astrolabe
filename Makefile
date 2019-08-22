all: build

build: arachne ivd kubernetes s3repository fs server cmd

cmd:
	cd cmd/arachne_server; go get; go install

arachne:
	cd pkg/arachne; go get; go install

ivd:
	cd pkg/ivd; go get; go install

fs:
	cd pkg/ivd; go get; go install

s3repository:
	cd pkg/s3repository; go get; go install

kubernetes:
	cd pkg/kubernetes; go get; go install

server:
	cd pkg/server; go get; go install
