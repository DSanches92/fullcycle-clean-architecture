# Clean Architecture | Listagem de Orders (REST, gRPC e GraphQL) - Go Expert

Neste desafio, você deve implementar a funcionalidade de **Listagem de Orders** em sua aplicação **Clean Architecture**. O objetivo principal é provar o desacoplamento da arquitetura: você criará um único Use Case (ListOrders) e o exporá através de três interfaces de comunicação diferentes simultaneamente.

## Tecnologias e Padrões

- **Linguagem:** Go (Golang)
- **Arquitetura:** Clean Architecture
- **Comunicação:** REST, gRPC e GraphQL
- **Infraestrutura:** Docker e Docker Compose

## Pré-requisitos

- Docker e Docker Compose instalados.
- Portas disponíveis na máquina:
  - REST: `3000`
  - GraphQL: `3001`
  - gRPC: `3002`

## Como Rodar

```bash
docker compose up --build
```

O comando sobe o banco (MySQL 8), aplica as migrações automaticamente e inicia a aplicação (que aguarda o banco estar pronto antes de rodar as migrações/servidores).

Portas dos serviços:

- **REST (Chi):** `http://localhost:3000`
  - `POST /orders` — criar pedido
  - `GET /orders` — listar pedidos
  - `GET /healthz` — liveness
  - `GET /readyz` — readiness
- **GraphQL (gqlgen):** `http://localhost:3001`
  - Playground: `http://localhost:3001/`
- **gRPC:** `localhost:3002`

Para testar, use as requisições prontas do arquivo `api.http` (criar e listar pedidos).

## Requisitos Técnicos

- [x] **Use Case:** Crie o caso de uso de listagem de pedidos (ListOrdersUseCase).
- [x] **Interfaces de Entrada:** Disponibilize o acesso a esse Use Case através de:
  - [x] REST: Endpoint GET /orders.
  - [x] gRPC: Service ListOrders.
  - [x] GraphQL: Query ListOrders.

## Banco de Dados

- [x] Crie as **migrações** necessárias para criar as tabelas do banco de dados.
- [x] O banco deve ser provisionado via Docker.

## Requisitos de Dockerização (Automação Total)

O avaliador não deve executar nenhum comando manual além do Docker Compose up, então:

- [x] Container da Aplicação: Você deve criar um Dockerfile para a sua aplicação Go.
- [x] Orquestração: O docker-compose.yaml deve subir o banco de dados e o container da aplicação.
- [x] Execução Automática: Ao rodar o comando docker compose up:
  - [x] O banco de dados deve subir.
  - [x] As migrações devem ser aplicadas automaticamente.
  - [x] A aplicação deve iniciar e ficar disponível nas portas configuradas.
  - [x] *Atenção:* Garanta que a aplicação aguarde o banco estar pronto antes de tentar rodar as migrações ou iniciar (handling de race condition na inicialização).

# Arquivos Auxiliares

- [x] Crie um arquivo `api.http` na raiz contendo as requisições prontas para:
  - [x] Criar uma Order (para popular o banco e testar).
  - [x] Listar as Orders (para validar o desafio).
