FROM golang:1.12-alpine3.9

COPY ./ /DynamicProxyEip/
RUN go get -v -u github.com/golang/dep/cmd/dep

WORKDIR /DynamicProxyEip
RUN dep ensure -v

RUN go build -ldflags '-w -s' -o dynamicProxyEip

FROM alpine:3.9

COPY --from=0 /DynamicProxyEip/dynamicProxyEip /

ENTRYPOINT ["/dynamicProxyEip"]