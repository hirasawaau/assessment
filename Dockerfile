FROM golang:1.19.4-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /main

ENV PORT 2565

CMD [ "./main" ]