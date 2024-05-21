FROM golang:1.21-alpine

# Instalar git
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY index.html ./
RUN go build -o /servicea

#FROM gcr.io/distroless/base-debian10
#COPY --from=builder /servicea /servicea
#COPY --from=builder /app/index.html /index.html

EXPOSE 8080

ENTRYPOINT ["/servicea"]
