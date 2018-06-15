FROM golang:1.10 as builder

RUN go get github.com/uphy/go-wso2am
WORKDIR /go/src/github.com/uphy/go-wso2am/wso2am-cli
RUN CGO_ENABLED=0 go build -o /wso2am-cli

FROM alpine:3.7

COPY --from=builder /wso2am-cli /bin/wso2am-cli
ENTRYPOINT [ "wso2am-cli" ]
