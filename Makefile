SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=2fanginx

#LDFLAGS=-ldflags "-X github.com/lstep/2fanginx/server.buildst=`date -u '+%Y-%m-%dT%I:%M:%S'` -X github.com/lstep/2fanginx/server.githash=`git rev-parse HEAD`"
LDFLAGS=-ldflags "-X github.com/lstep/2fanginx/server.Build=`git rev-parse HEAD`"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

test:
	go test

