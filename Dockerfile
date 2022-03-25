FROM golang:1.17-alpine as builder
RUN apk add git openssh-client
WORKDIR /go/src/app
COPY . .
RUN mkdir -p /root/.ssh \
    && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

RUN go mod tidy
RUN go build -o app

FROM alpine:latest
RUN apk add tzdata
RUN apk add --no-cache ca-certificates
RUN adduser -D niko
USER niko
WORKDIR /go/src/app
COPY --from=builder /go/src/app/app /usr/local/bin/


EXPOSE 9180
ENTRYPOINT ["app"]
