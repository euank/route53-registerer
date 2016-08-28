.PHONY: all static docker push clean

all:
	go build .

static: bin/route53-registerer

bin/route53-registerer:
	CGO_ENABLED=0 go build -installsuffix cgo -a -ldflags "-s" -o ./bin/route53-registerer


bin/cacert.pem:
	mkdir -p bin
	wget -O bin/cacert.pem https://curl.haxx.se/ca/cacert.pem

docker: bin/cacert.pem
	docker build -t euank/route53-registerer-builder:nopush -f Dockerfile.build .
	docker run -v "$(PWD)/bin":/go/src/github.com/euank/route53-registerer/bin euank/route53-registerer-builder:nopush
	docker build -t euank/route53-registerer:latest -f Dockerfile.release .
	docker tag euank/route53-registerer:latest euank/route53-registerer:$(shell cat VERSION)

push:
	docker push euank/route53-registerer:$(shell cat VERSION)
	docker push euank/route53-registerer:latest 

clean:
	rm -f bin/cacert.pem bin/route53-registerer
	docker rmi euank/route53-registerer
