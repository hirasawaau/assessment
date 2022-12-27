FROM golang:1.19.3-alpine as build

WORKDIR /go/src/github.com/hirasawaau/assessment

COPY go.* ./
COPY *.go ./
COPY src ./src
ENV PORT=2565
RUN go mod download
RUN CGO_ENABLED=0 go build -o a.out server.go


FROM alpine:3.14.2
WORKDIR /app
COPY --from=build /go/src/github.com/hirasawaau/assessment/a.out .
ENV PORT=2565
EXPOSE 2565
CMD ["./a.out"]