FROM golang:1.25.6-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o booking-service

EXPOSE 50051

CMD [ "./booking-service" ]