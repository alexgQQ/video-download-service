FROM golang:alpine as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

# Provides the compression of a scratch image
# but includes the needed tools (openssl) for
# various actions
# https://github.com/GoogleContainerTools/distroless/blob/main/base/README.md
FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=builder /app/main /usr/bin/
COPY templates templates
COPY public public

ENTRYPOINT ["main"]