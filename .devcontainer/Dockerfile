FROM golang:1.11.0-stretch  as builder
RUN mkdir /build
ADD ./cmd/main/ /build/main/
ADD ./cmd/elasticsearch/ /build/elasticsearch/
ADD ./cmd/filterlines/ /build/filterlines/
ADD ./cmd/parsecontents/ /build/parsecontents/
ADD ./cmd/parsedates/ /build/parsedates/
WORKDIR /build
COPY go.mod go.sum ./
WORKDIR /build/main
RUN GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o server.exe .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/server /app/
WORKDIR /app
CMD ["./server"]