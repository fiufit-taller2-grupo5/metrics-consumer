FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

ENV environment=production

RUN go build -o /main .
