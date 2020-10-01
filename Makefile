GOPKG ?=	moul.io/zapconfig

include rules.mk

generate:
	GO111MODULE=off go get github.com/campoy/embedmd
	embedmd -w README.md
