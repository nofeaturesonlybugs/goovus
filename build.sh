#!/bin/bash
DATE=$(date)
GOVERSION=$(go version)
VERSION=$(git describe --tags --abbrev=8 --dirty --always --long)

LDFLAGS=
LDFLAGS="$LDFLAGS -X 'main.Version=$VERSION'"
LDFLAGS="$LDFLAGS -X 'main.Date=$DATE'"
LDFLAGS="$LDFLAGS -X 'main.GoVersion=$GOVERSION'"
go build -ldflags "$LDFLAGS"
