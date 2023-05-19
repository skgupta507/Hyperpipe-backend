FROM --platform=$BUILDPLATFORM docker.io/golang:latest AS build

WORKDIR /src

RUN apk --no-cache add ca-certificates

COPY . .

ARG TARGETOS TARGETARCH

RUN GOPROXY=https://proxy.golang.org,https://goproxy.cn go mod download && \
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o /src/hyperpipe-backend

FROM scratch as bin

WORKDIR /app
COPY --from=build /etc/ssl/certs/cacertificates.crt /etc/ssl/certs/
COPY --from=build /src/hyperpipe-backend .

EXPOSE 3000

CMD ./hyperpipe-backend
