FROM golang:1.21-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

WORKDIR /app

RUN go build -o /app/cmd .

FROM golang:1.21-alpine

WORKDIR /app

COPY --from=build /app/cmd .

ENTRYPOINT ["./cmd"]