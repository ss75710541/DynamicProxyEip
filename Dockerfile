FROM golang:1.12

COPY ./ /DynamicProxyEip/
RUN  go get -v -u github.com/golang/dep/cmd/dep

WORKDIR /DynamicProxyEip
RUN dep ensure -v

RUN GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o dynamicProxyEip

FROM alpine:3.9

COPY --from=0 /DynamicProxyEip/dynamicProxyEip /

ENTRYPOINT ["/dynamicProxyEip"]