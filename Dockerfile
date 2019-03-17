FROM golang:1.11.5-alpine as builder
ENV CGO_ENABLED=0
RUN apk add --no-cache git
RUN mkdir -p /usr/local/app
WORKDIR /usr/local/app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . /usr/local/app
RUN go test -v ./...
RUN go install -v ./cmd/...

FROM alpine
RUN apk add --no-cache ca-certificates tzdata
ENV ADDR=":80"
EXPOSE 80
RUN mkdir -p /usr/local/app
WORKDIR /usr/local/app
ENTRYPOINT ["./dflimg"]
COPY --from=builder /go/bin/dflimg .
