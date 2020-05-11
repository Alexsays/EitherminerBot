FROM golang:1.14-alpine3.11

RUN apk add --no-cache git

WORKDIR /app/eithermine

COPY go.mod .

RUN go mod download

COPY . .

ENTRYPOINT ["./wait-for.sh", "postgres:5432", "--", "go", "run", "main.go"]
