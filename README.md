# GoExpert Stress Test

Projeto CLI em Go para executar testes de carga (stress tests) em serviços HTTP, enviando um número configurável de requisições de forma concorrente e exibindo um relatório consolidado no console.

## Funcionalidades

- Disparo de requisições HTTP GET com quantidade configurável
- Controle de concorrência via flag `--concurrency`
- Relatório contendo:
  - Tempo total de execução
  - Total de requisições realizadas
  - Quantidade de respostas `HTTP 200`
  - Distribuição completa dos status codes retornados
- Execução simplificada via Docker ou binário Go
- Testes unitários prontos (`go test`)

## Configuração

A ferramenta é configurada **apenas** por flags de linha de comando:

| Flag | Valor padrão | Descrição |
|------|--------------|-----------|
| `--url` | *obrigatório* | URL do serviço a ser testado |
| `--requests` | `1` | Número total de requisições a serem enviadas |
| `--concurrency` | `1` | Número de requisições simultâneas |

Observações:

- `--concurrency` não pode exceder `--requests`; se exceder, será ajustado automaticamente.
- Somente requisições `GET` são disparadas.

## Como rodar localmente (sem Docker)

Pré-requisito: Go 1.22+ instalado.

```bash
go run ./cmd/stresstest --url http://google.com --requests 1000 --concurrency 10
```

## Como rodar com Docker

```bash
# Build da imagem
docker build -t stresstest .

# Execução
docker run --rm stresstest --url http://google.com --requests 1000 --concurrency 10
```

## Exemplo de saída

```console
===== Relatório =====
URL: http://google.com
Tempo total: 2.143127s
Total de requests: 1000
Requests com HTTP 200: 1000
Distribuição de status codes:
  200 -> 1000
```

## Testes unitários

```bash
go test ./...
```
