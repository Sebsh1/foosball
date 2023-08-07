FROM golang:1.20-alpine3.18 AS builder
WORKDIR /build
COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/bin/service main.go

FROM alpine:3.18 as image-base
WORKDIR /app
COPY --from=builder /build/bin/service /usr/bin/service
ENTRYPOINT [ "service" ]
CMD [ "serve" ]