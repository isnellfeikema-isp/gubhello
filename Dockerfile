FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/isnellfeikema-isp/gubhello/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ["go", "install", "."]

FROM alpine
COPY --from=builder /go/bin/gubhello /usr/local/bin/

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/gubhello"]