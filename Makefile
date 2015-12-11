all:
	go build -ldflags "-X server.buildst=`date -u '+%Y-%m-%dT%I:%M:%S'`"

#go build -ldflags "-X main.buildstamp0=`date -u '+%Y-%m-%dT%I:%M:%S'` -X main.githash0=`git rev-parse HEAD`"

test:
	go test

