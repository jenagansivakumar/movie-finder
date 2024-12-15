FROM golang:1.20-alipine

WORKDIR /app

COPY ./ /app/

RUN  go mod tidy