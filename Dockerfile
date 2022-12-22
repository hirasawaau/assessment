FROM golang:1.19.3-alpine

WORKDIR /go/src/github.com/hirasawaau/assessment

COPY go.* ./
COPY *.go ./
COPY src ./src
ENV PORT=2565
RUN go mod download
RUN CGO_ENABLED=0 go build -o a.out server.go


EXPOSE 2565

CMD ["./a.out"]