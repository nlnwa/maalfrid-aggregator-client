FROM golang:alpine

RUN apk add --no-cache --update alpine-sdk

COPY . /go/src/github.com/nlnwa/maalfrid-aggregator-client

RUN cd /go/src/github.com/nlnwa/maalfrid-aggregator-client \
&& go get ./... \
&& VERSION=$(./scripts/git-version) \
CGO_ENABLED=0 \
go install -a -tags netgo -v -ldflags "-w -X github.com/nlnwa/maalfrid-aggregator-client/pkg/version.Version=$VERSION" \
github.com/nlnwa/maalfrid-aggregator-client/cmd/...
# -w Omit the DWARF symbol table.
# -X Set the value of the string variable in importpath named name to value.


FROM scratch

LABEL maintainer="marius.beck@nb.no"

COPY --from=0 /go/bin/maalfrid-aggregator-client /

ENV HOST=localhost \
    PORT=8672

ENTRYPOINT ["/maalfrid-aggregator-client"]
CMD ["version"]
