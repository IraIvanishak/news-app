FROM golang:1.25.3 AS builder

WORKDIR /go/src/app

COPY . .

RUN go mod download && \
    go build -v -o api cmd/api/main.go 

FROM builder AS production

WORKDIR /bin

COPY --from=builder /go/src/app/api .
COPY --from=builder /go/src/app/web ./web
COPY --from=builder /go/src/app/.env .

CMD api

FROM builder AS dev
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
WORKDIR /go/src/app

CMD air
