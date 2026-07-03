# Full Cycle - Pós Go Expert - Client/Server API

Projeto desenvolvido em Go com foco na comunicação entre cliente e servidor utilizando HTTP, Context, SQLite e manipulação de arquivos.

O objetivo deste projeto é demonstrar:

-   Criação de servidor HTTP
-   Consumo de API externa
-   Utilização de Context para controle de timeout
-   Persistência de dados utilizando SQLite
-   Comunicação entre aplicações Client/Server
-   Manipulação de arquivos em Go

------------------------------------------------------------------------

# Tecnologias utilizadas

-   Go 1.25
-   HTTP Server
-   HTTP Client
-   Context
-   SQLite
-   JSON
-   Manipulação de Arquivos

------------------------------------------------------------------------

# Estrutura do projeto

``` text
.
├── client.go
├── server
│   └── server.go
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

------------------------------------------------------------------------

# Funcionamento da aplicação

O projeto é composto por duas aplicações independentes.

## Server

O servidor disponibiliza o endpoint:

``` text
GET /cotacao
```

Ao receber uma requisição:

1.  Consome a API de câmbio AwesomeAPI.
2.  Limita a chamada utilizando Context com timeout de **200ms**.
3.  Persiste a cotação em um banco SQLite.
4.  Limita a gravação no banco utilizando Context com timeout de
    **10ms**.
5.  Retorna ao cliente apenas o valor do campo **bid** em formato JSON.

Exemplo de resposta:

``` json
{
    "bid": "5.4321"
}
```

------------------------------------------------------------------------

## Client

O cliente realiza uma requisição ao servidor local utilizando Context
com timeout de **300ms**.

Após receber a resposta:

-   realiza o parse do JSON;
-   recupera apenas o campo **bid**;
-   grava o resultado no arquivo:

``` text
cotacao.txt
```

Formato do arquivo:

``` text
Dólar: 5.4321
```

------------------------------------------------------------------------

# Timeouts

| Operação                      | Timeout         |
| ----------                    | --------------- |
| Consumo da API externa        |          200 ms |
| Persistência no SQLite        |           10 ms |
| Requisição do Client          |          300 ms |

Quando um timeout é excedido, uma mensagem é registrada no console da aplicação correspondente.

------------------------------------------------------------------------

# Como executar

## 1. Clonar o repositório

``` bash
git clone git@github.com:danielsantello/fc-pos-desafio-client-server-api.git
```

Entre na pasta do projeto:

``` bash
cd fc-pos-desafio-client-server-api
```

## 2. Executar o servidor

``` bash
go run server/server.go
```

Servidor disponível em:

``` text
http://localhost:8080/cotacao
```

## 3. Executar o cliente

``` bash
go run client.go
```

Após a execução será criado o arquivo:

``` text
cotacao.txt
```

------------------------------------------------------------------------

# Banco de dados

O banco SQLite é criado automaticamente durante a inicialização do
servidor.

Tabela utilizada:

``` sql
CREATE TABLE cotacao (
    id INTEGER PRIMARY KEY,
    bid TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

Cada requisição realizada ao endpoint gera um novo registro na tabela.

------------------------------------------------------------------------

# Arquivos gerados

Durante a execução da aplicação são criados automaticamente:

``` text
cotacao.db
cotacao.txt
```

Esses arquivos representam artefatos gerados pela execução e, por esse
motivo, não fazem parte do repositório, estando configurados no arquivo
`.gitignore`.

------------------------------------------------------------------------

# Autor

Daniel Santello
