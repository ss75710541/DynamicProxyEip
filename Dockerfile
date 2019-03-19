FROM golang:1.12

RUN  go get -v -u github.com/golang/dep/cmd/dep

COPY ./ /go/src/github.com/DynamicProxyEip
WORKDIR /go/src/github.com/DynamicProxyEip

RUN dep ensure -v

RUN GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o dynamicProxyEip

FROM alpine:3.9

COPY --from=0 /go/src/github.com/DynamicProxyEip/dynamicProxyEip /

ENTRYPOINT ["/dynamicProxyEip"]
