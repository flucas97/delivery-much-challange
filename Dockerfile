FROM golang:1.13-buster AS builder

LABEL maintainer="Lucas Feijo <feijolucas1997@gmail.com>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api -ldflags "-s -w" -a -installsuffix cgo cmd/main.go

FROM debian:buster
RUN apt update -y && apt upgrade -y && apt install ca-certificates -y
WORKDIR /app

COPY --from=builder /app/entrypoint.sh /app/
COPY --from=builder /app/api /app/

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]