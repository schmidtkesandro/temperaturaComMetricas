# Sistema de Temperatura por CEP

## Descrição

Este projeto consiste em dois serviços que se comunicam para fornecer a temperatura de uma cidade com base em um CEP fornecido.

- **Serviço A:** Recebe o CEP e encaminha para o Serviço B.
- **Serviço B:** Pesquisa o CEP, obtém a cidade e retorna a temperatura.

## Como Rodar

### Requisitos

- Docker
- Docker Compose

### Passos

1. Clone o repositório.
2. Navegue até o diretório do projeto.
3. Execute o comando `docker-compose up --build`.
4. O Serviço A estará disponível em `http://localhost:8080` e o Serviço B em `http://localhost:8081`.
5. Para parar os containers 'docker-compose down'.

### Exemplos de Uso

- **Serviço A:**
  - Endpoint: `/cep`
  - Método: `POST`
  - Body: `{ "cep": "29902555" }`

curl -X POST <http://localhost:8080/cep> -H "Content-Type: application/json" -d '{"cep": "29902555"}'

### Exemplos de Uso no Browser

- **Busca por CEP:**
  - <http://localhost:8080/>

- **BPara ver as métricas:**
  - <http://localhost:8080/metrics>

### OpenTelemetry e Zipkin

O projeto inclui tracing distribuído com OpenTelemetry e Zipkin. Para visualizar os traces, acesse o Zipkin em <http://localhost:9411>.

### Grafana e Prometheus

Para visualizar o grafana em <http://localhost:3000> e o Prometheus em <http://localhost:9090>, Usuário e senha no Prometheus é admin/admin.
