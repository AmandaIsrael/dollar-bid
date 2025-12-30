# Dollar Bid - Sistema de Cotação do Dólar

Este projeto implementa dois sistemas separados em Go para consulta e registro de cotação do dólar.

## Estrutura do Projeto

```
dollar-bid/
├── server/           # Servidor HTTP que busca cotações
│   ├── main.go
│   ├── server.go
│   └── go.mod
├── client/           # Cliente que consulta o servidor
│   ├── main.go
│   └── go.mod
└── README.md
```

## Como Executar

### 1. Iniciando o Servidor

```bash
cd server
go mod tidy
go run .
```

O servidor será iniciado na porta 8080 e exibirá:
```
Servidor iniciado na porta 8080
```

### 2. Executando o Cliente

Em outro terminal:

```bash
cd client
go run .
```

O cliente fará a requisição ao servidor e salvará a cotação em `cotacao.txt`.

## Funcionalidades

### Servidor (`server/`)
- **Endpoint**: `GET http://localhost:8080/cotacao`
- **API Externa**: `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- **Timeout API externa**: 200ms
- **Timeout banco de dados**: 10ms
- **Banco**: SQLite (`cotacoes.db`)
- **Logs**: Registra todas as operações e timeouts

### Cliente (`client/`)
- **Timeout requisição**: 300ms
- **Arquivo de saída**: `cotacao.txt`
- **Formato**: `Dólar: {valor}`
- **Logs**: Registra operações e erros

## Timeouts Implementados

1. **Cliente → Servidor**: 300ms
2. **Servidor → API Externa**: 200ms  
3. **Servidor → Banco de Dados**: 10ms

Todos os timeouts geram logs específicos quando excedidos.

## Exemplo de Uso

1. Execute o servidor
2. Execute o cliente
3. Verifique o arquivo `cotacao.txt` criado:

```
Dólar: 6.1234
```

## Logs de Exemplo

**Servidor:**
```
[DATABASE] Conexão com SQLite estabelecida com sucesso
[SERVER] Cotação obtida da API externa com sucesso
[DATABASE] Cotação salva no banco com sucesso. ID: 1
[SERVER] Cotação retornada com sucesso
```

**Cliente:**
```
[CLIENT] Cotação salva no arquivo cotacao.txt: 6.1234
```