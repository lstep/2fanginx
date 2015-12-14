all:
	go build -ldflags "-X 2fanginx/server.buildst=`date -u '+%Y-%m-%dT%I:%M:%S'` -X 2fanginx/server.githash=`git rev-parse HEAD`"

test:
	go test

