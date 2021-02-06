FROM golang:1.15-alpine
RUN apk add --no-cache git
RUN apk add g++

WORKDIR /tradeX

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GO_ENABLED=1
RUN go install github.com/mattn/go-sqlite3


RUN go build -o ./out/tradeX ./main

EXPOSE 1778

CMD ["./out/tradeX"]