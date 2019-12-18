FROM golang:alpine as builder
RUN mkdir -p /build/release
ADD . /build
WORKDIR /build
RUN go build -o /build/release ./...
FROM alpine
COPY --from=builder /build/release/gtt /usr/local/bin/gtt
ENTRYPOINT ["/usr/local/bin/gtt"]
