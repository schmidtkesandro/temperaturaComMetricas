FROM golang:1.21-alpine

# Instalar git
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./


RUN go build -o /serviceb

#FROM gcr.io/distroless/base-debian10
#COPY --from=builder /serviceb /serviceb

EXPOSE 8081

ENTRYPOINT ["/serviceb"]
