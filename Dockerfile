FROM golang:alpine AS build

ARG proxy

WORKDIR /app/

RUN apk --no-cache add ca-certificates

COPY . .

RUN sed -i "s/hyperpipe-proxy.onrender.com/$proxy/g" utils.go

RUN go mod download && \
	go build -ldflags "-s -w"

EXPOSE 3000

CMD ./hyperpipe-backend