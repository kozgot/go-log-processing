FROM golang:1.13.0-stretch  as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
COPY go.mod ./
WORKDIR /app
RUN GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /app/server /app/
WORKDIR /app
CMD ["./server"] 