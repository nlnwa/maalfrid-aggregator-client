FROM golang:alpine

RUN apk add --no-cache --update alpine-sdk

COPY . /src/maalfrid-aggregator-client

# build flags:
#  -w Omit the DWARF symbol table.
#  -X Set the value of the string variable in importpath named name to value.

RUN cd /src/maalfrid-aggregator-client \
&& VERSION=$(./scripts/git-version) \
CGO_ENABLED=0 \
go install \
-a -tags netgo -v -ldflags "-w -X github.com/nlnwa/maalfrid-aggregator-client/pkg/version.Version=$VERSION" \
github.com/nlnwa/maalfrid-aggregator-client/cmd/...


FROM scratch

LABEL maintainer="marius.beck@nb.no"

COPY --from=0 /go/bin/maalfrid-aggregator-client /

ENV HOST=localhost \
    PORT=8672

ENTRYPOINT ["/maalfrid-aggregator-client"]
CMD ["version"]
