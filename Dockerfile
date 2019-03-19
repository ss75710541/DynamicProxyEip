FROM golang:1.12

RUN  go get -v -u github.com/golang/dep/cmd/dep

COPY ./ /go/src/github.com/DynamicProxyEip
WORKDIR /go/src/github.com/DynamicProxyEip

RUN dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o dynamicProxyEip

FROM alpine:3.9

RUN apk add --no-cache ca-certificates

COPY --from=0 /go/src/github.com/DynamicProxyEip/dynamicProxyEip /

ENTRYPOINT ["/dynamicProxyEip"]
