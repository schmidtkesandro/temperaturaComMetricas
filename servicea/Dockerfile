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

EXPOSE 8080

CMD [ "/servicea" ]
