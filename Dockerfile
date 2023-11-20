FROM golang:1.21 AS builder
ENV CGO_ENABLED=0

WORKDIR /src
COPY . /src

RUN go build -o "./bin/server" "./app/server"

# We don't need golang to run binaries, just use alpine.
FROM alpine
COPY --from=builder /src/bin/server /app/server
COPY --from=builder /src/start.sh /app/start.sh

EXPOSE 8080

WORKDIR /app

CMD sh start.sh
