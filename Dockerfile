FROM golang:alpine

ENV GO111MODULE=on

RUN apk add --no-cache --update alpine-sdk
WORKDIR /go/src/github.com/nlnwa/maalfrid-aggregator-client

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# build flags:
#  -w Omit the DWARF symbol table.
#  -X Set the value of the string variable in importpath named name to value.
RUN VERSION=$(./scripts/git-version) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -v \
-ldflags "-w -X github.com/nlnwa/maalfrid-aggregator-client/pkg/version.Version=${VERSION}"


FROM scratch

LABEL maintainer="marius.beck@nb.no"

COPY --from=0 /go/bin/maalfrid-aggregator-client /

ENV HOST=localhost \
    PORT=8672

ENTRYPOINT ["/maalfrid-aggregator-client"]
CMD ["version"]
