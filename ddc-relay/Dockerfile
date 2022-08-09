FROM golang:1.17 AS builder
ENV CGO_ENABLED=1 GOOS=linux
WORKDIR /go/src/ddc-relayy
RUN apt-get update -y && apt-get install ca-certificates -y
RUN go get github.com/bugsnag/panic-monitor
COPY . .
RUN go build -a -o /go/bin/ddc-relayy ./cmd/ddc-relayy

FROM alpine:latest
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add libc6-compat curl
COPY --from=builder /go/bin/ddc-relayy /go/bin/ddc-relayy
COPY --from=builder /go/bin/panic-monitor /go/bin/panic-monitor
ENTRYPOINT ["/go/bin/panic-monitor", "/go/bin/ddc-relayy"]
CMD ["run"]
