# Define a imagem base para o estágio de compilação
FROM golang:1.21 AS builder

# Define o diretório de trabalho no container
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia os arquivos do código fonte
COPY . .

# Compila o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o servicea

# Define a imagem base para o estágio de execução
FROM gcr.io/distroless/base-debian10

# Copia o executável do estágio de compilação
COPY --from=builder /app/servicea /

# Define o comando para executar o aplicativo
CMD ["/servicea"]
