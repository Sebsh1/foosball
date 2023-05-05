FROM golang:1.18-alpine3.17 AS builder
WORKDIR /build
COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/bin/service main.go

FROM alpine:3.17 as image-base
WORKDIR /app
COPY --from=builder /build/bin/service /usr/bin/service
ENTRYPOINT [ "service" ]
CMD [ "serve" ]