FROM golang:1.21-alpine as build
WORKDIR /workspace
COPY . /workspace
RUN go mod download
RUN apk add build-base
RUN go build -ldflags="-w -s" -o /thumburl-service ./cmd

FROM alpine as thumburl-service
COPY --from=build /thumburl-service /thumburl-service
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/thumburl-service"]
