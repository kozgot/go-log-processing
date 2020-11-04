FROM golang:1.13.0-stretch  as builder
RUN mkdir /build
ADD ./cmd/main/ /build/cmd/main/
ADD ./cmd/elasticsearch/ /build/cmd/elasticsearch/
ADD ./cmd/filterlines/ /build/cmd/filterlines/
ADD ./cmd/parsecontents/ /build/cmd/parsecontents/
ADD ./cmd/parsedates/ /build/cmd/parsedates/
WORKDIR /build
COPY go.mod go.sum ./
WORKDIR /build/cmd/main
RUN GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY ./cmd/main/*.log /app/
COPY --from=builder /build/cmd/main/server /app/
WORKDIR /app
ENTRYPOINT ["./server"]