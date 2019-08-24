all: build

build: arachne ivd kubernetes s3repository fs server cmd

deps:
	go get k8s.io/klog
	go get k8s.io/api/core/v1
	go get k8s.io/apimachinery/pkg/apis/meta/v1
	go get k8s.io/client-go/tools/clientcmd
	go get k8s.io/client-go/kubernetes
	go get github.com/aws/aws-sdk-go
	go get github.com/pkg/errors
	go get github.com/vmware/govmomi
	go get github.com/google/uuid
	go get github.com/labstack/echo
	cd $(GOPATH)/src/k8s.io/klog ; git checkout v0.4.0

cmd: deps
	cd cmd/arachne_server; go build

arachne: deps
	cd pkg/arachne; go build

ivd: deps
	cd pkg/ivd; go build

fs: deps
	cd pkg/ivd; go build

s3repository: deps
	cd pkg/s3repository; go build

kubernetes: deps
	cd pkg/kubernetes; go build

server: deps
	cd pkg/server; go build
