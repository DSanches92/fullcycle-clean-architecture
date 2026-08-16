# Clean Architecture | Listagem de Orders (REST, gRPC e GraphQL) - Go Expert

Neste desafio, você deve implementar a funcionalidade de **Listagem de Orders** em sua aplicação **Clean Architecture**. O objetivo principal é provar o desacoplamento da arquitetura: você criará um único Use Case (ListOrders) e o exporá através de três interfaces de comunicação diferentes simultaneamente.

## Tecnologias e Padrões

- **Linguagem:** Go (Golang)
- **Arquitetura:** Clean Architecture
- **Comunicação:** REST, gRPC e GraphQL
- **Infraestrutura:** Docker e Docker Compose

## Pré-requisitos

TODO: Adicionar aqui os pre requisitos.

## Como Rodar

> TODO: Adicionar aqui, como rodar o projeto.
> - O comando único de execução (docker compose up).
> - As portas em que cada serviço (Web, gRPC, GraphQL) está rodando.

## Requisitos Técnicos

- [ ] **Use Case:** Crie o caso de uso de listagem de pedidos (ListOrdersUseCase).
- [ ] **Interfaces de Entrada:** Disponibilize o acesso a esse Use Case através de:
  - [ ] REST: Endpoint GET /orders.
  - [ ] gRPC: Service ListOrders.
  - [ ] GraphQL: Query ListOrders.

## Banco de Dados

- [ ] Crie as **migrações** necessárias para criar as tabelas do banco de dados.
- [ ] O banco deve ser provisionado via Docker.

## Requisitos de Dockerização (Automação Total)

O avaliador não deve executar nenhum comando manual além do Docker Compose up, então:

- [ ] Container da Aplicação: Você deve criar um Dockerfile para a sua aplicação Go.
- [ ] Orquestração: O docker-compose.yaml deve subir o banco de dados e o container da aplicação.
- [ ] Execução Automática: Ao rodar o comando docker compose up:
  - [ ] O banco de dados deve subir.
  - [ ] As migrações devem ser aplicadas automaticamente.
  - [ ] A aplicação deve iniciar e ficar disponível nas portas configuradas.
  - [ ] *Atenção:* Garanta que a aplicação aguarde o banco estar pronto antes de tentar rodar as migrações ou iniciar (handling de race condition na inicialização).

# Arquivos Auxiliares

- [ ] Crie um arquivo `api.http` na raiz contendo as requisições prontas para:
  - [ ] Criar uma Order (para popular o banco e testar).
  - [ ] Listar as Orders (para validar o desafio).
