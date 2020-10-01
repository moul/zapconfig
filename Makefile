GOPKG ?=	moul.io/zapconfig
DOCKER_IMAGE ?=	moul/zapconfig
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate: install
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	echo 'foo@bar:~$$ zapconfig' > .tmp/usage.txt
	zapconfig 2>&1 >> .tmp/usage.txt
	embedmd -w README.md
	rm -rf .tmp
