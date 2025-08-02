# 🚀 labs-auction-goexpert

Sistema de **leilões em Go** (desafio FullCycle GoExpert), com **fechamento automático** do leilão após tempo configurável via variável de ambiente.

---

## 🧩 Visão geral

O projeto implementa APIs REST para:

- Criar leilões (auction)
- Realizar lances (bids)
- Consultar leilões e lances
- Obter lance vencedor
- Fechamento automático dos leilões após expiração

As principais regras de negócio estão nas camadas `internal/usecase` e `internal/repository`. O fechamento automático se baseia na criação de uma **goroutine** acionada no momento da criação do leilão, que aguarda o tempo — definido por `AUCTION_DURATION_SECONDS` — e então fecha o leilão, impedindo novos lances.

O projeto foi estruturado com **Docker e Docker Compose**, incluindo um container MongoDB para persistência :contentReference[oaicite:1]{index=1}.

---

## Como Rodar o Projeto
1. Clone o repositório
```
git clone https://github.com/seu-usuario/labs-auction-goexpert.git
cd labs-auction-goexpert
```
2. Construa e inicie os containers
```
docker compose up -d --build
docker compose up
```

📁 Estrutura de Diretórios
```
.
├── cmd
│   └── auction
│       └── main.go                # Ponto de entrada principal da aplicação (inicializa o servidor HTTP)
├── configuration
│   ├── database
│   │   └── mongodb                # Configuração de conexão com MongoDB
│   ├── logger                     # Configuração do logger padrão do sistema
│   └── rest_err                   # Tratamento centralizado de erros REST (HTTP)
├── docker-compose.yml            # Subida dos containers (app + MongoDB)
├── Dockerfile                    # Dockerfile para build da aplicação Go
├── go.mod                        # Gerenciador de dependências (Go Modules)
├── go.sum
├── internal
│   ├── entity
│   │   ├── auction_entity         # Entidade de domínio Auction
│   │   ├── bid_entity             # Entidade de domínio Bid
│   │   └── user_entity            # Entidade de domínio User
│   ├── infra
│   │   ├── api
│   │   │   └── web
│   │   │       └── controller     # Camada de entrada da aplicação (handlers HTTP via Gin)
│   │   ├── database
│   │   │   ├── auction            # Implementações de persistência da entidade Auction (Mongo)
│   │   │   ├── bid                # Implementações de persistência da entidade Bid
│   │   │   └── user               # Implementações de persistência da entidade User
│   ├── internal_error             # Definições de erros internos da aplicação (usado por toda a lógica)
│   └── usecase
│       ├── auction_usecase        # Casos de uso relacionados a leilões
│       ├── bid_usecase            # Casos de uso relacionados a lances
│       └── user_usecase           # Casos de uso relacionados a usuários
```

🔧 Variáveis de ambiente
AUCTION_DURATION_SECONDS: tempo em segundos para o leilão permanecer aberto após criação; se ausente ou inválida, o default de 30s é utilizado.

Além disso, pode-se configurar PORT (porta HTTP, padrão 8080) ou MONGO_URI dependendo da estrutura do projeto.

📡 Endpoints da API
Todas as rotas usando Content-Type: application/json. O servidor Gin (arquivo cmd/auction/main.go) define nosso router principal.

1. Criar um leilão
```
POST /auctions
```

Request Body:
```
{
  "product_name": "MacBook Pro 2019",
  "category": "computers",
  "description": "i9, 16GB RAM, 512GB SSD",
  "condition": "new"
}
```
Response 201 Created:
```
{
  "id": "uuid",
  "product_name": "...",
  "category": "...",
  "description": "...",
  "condition": "new",
  "status": "open",
  "timestamp": "2025-08-21T13:45:30Z"
}
```

2. Listar todos os leilões
```
GET /auctions
```
Retorna array de objetos Auction.

3. Consultar um leilão por ID
```
GET /auctions/:id
```
4. Realizar um lance
```
POST /auctions/:id/bids
```
Request Body:
```
{
  "user_id": "uuid do usuário",
  "amount": 125.50
}
```
Respostas possíveis:
201 Created, se lance aceito (valor maior que o anterior e leilão ainda aberto).
400 Bad Request, caso lance inválido ou leilão fechado.

5. Listar lances de um leilão
```
GET /auctions/:id/bids
```
Retorna array de objetos Bid (com id, user_id, amount, timestamp).

6. Obter lance vencedor
```
GET /auctions/:id/bids/winner
```
Retorna o Bid com maior valor registrado (mesmo que leilão ainda esteja aberto).

7. criar usuario 
```
POST /user
```

Request Body:
```
{
    "name": "nome usuario"
}
```

## 📦 Tecnologias e requisitos

- Go 1.18+ (estrutura com modules compatível com Go 1.19+)
- Docker & Docker Compose
- MongoDB (embutido via `docker-compose.yml`)
- Dotenv (.env) contendo variável `AUCTION_DURATION_SECONDS` (default: 30)
- Testes unitários com Go, incluindo fluxos de criação e fechamento automático

---


