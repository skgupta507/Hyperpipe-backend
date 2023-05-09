FROM --platform=$BUILDPLATFORM golang:alpine AS build

WORKDIR /app/

RUN apk --no-cache add ca-certificates

COPY . .

ARG TARGETOS TARGETARCH

RUN GOPROXY=https://goproxy.cn go get github.com/klauspost/compress && \
	CGO=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w"

EXPOSE 3000

CMD ./hyperpipe-backend