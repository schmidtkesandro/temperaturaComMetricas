FROM golang:1.18-alpine

# Instalar git
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /servicea

EXPOSE 8080

CMD [ "/servicea" ]