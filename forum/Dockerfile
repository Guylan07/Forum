FROM golang:1.21

WORKDIR /app

RUN apt-get update && apt-get install -y gcc libsqlite3-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p certs static/uploads

RUN go build -o forum ./cmd/main.go

EXPOSE 8085 8443

CMD ["./forum"]