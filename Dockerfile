FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY endpoints /
WORKDIR /app
ENTRYPOINT ["/endpoints"]
