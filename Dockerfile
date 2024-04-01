FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

WORKDIR /app

RUN go build -o /app/cmd .

FROM golang:1.22-alpine

WORKDIR /app

COPY --from=build /app/cmd .

EXPOSE 8080

CMD ["./cmd"]